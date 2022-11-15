package pull

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	internalHttp "github.com/030/p2iwd/internal/pkg/http"
	archiverV4 "github.com/mholt/archiver/v4"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	blobs        = "/blobs/"
	uriCatalog   = "_catalog"
	uriManifests = "/manifests/"
	uriTagsList  = "/tags/list"
)

type DockerRegistry struct {
	Dir, Host, Pass, User string
}

type manifestJSON struct {
	Config   string   `json:"config"`
	RepoTags []string `json:"repoTags"`
	Layers   []string `json:"layers"`
}

func (dr *DockerRegistry) download(file, header, url string) error {
	if _, err := os.Stat(file); err == nil {
		log.Debugf("file: '%s' exists", file)
		return nil
	}

	ha := internalHttp.Auth{HeaderKey: "Accept", HeaderValue: header, Method: "GET", Pass: dr.Pass, User: dr.User, URL: url}
	rc, err := ha.RequestAndResponseBody(nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	f, err := os.Create(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	_, err = io.Copy(f, rc)
	if err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	log.Infof("file: '%s' created", file)

	return nil
}

func (dr *DockerRegistry) downloadLayer(blobSum, file, repo string) error {
	url := dr.Host + internalHttp.Version + repo + blobs + blobSum
	if err := dr.download(file, "application/vnd.docker.image.rootfs.diff.tar.gzip", url); err != nil {
		return err
	}

	return nil
}

func (dr *DockerRegistry) downloadManifest(dir, repo, tag string) error {
	file := filepath.Join(dir, "upload-manifest.json")
	url := dr.Host + internalHttp.Version + repo + uriManifests + tag
	header := "application/vnd.docker.distribution.manifest.v2+json"
	if err := dr.download(file, header, url); err != nil {
		return err
	}

	return nil
}

func (dr *DockerRegistry) downloadConfig(dir, repo, tag string) error {
	url := dr.Host + internalHttp.Version + repo + uriManifests + tag
	log.Info(url)
	header := "application/vnd.docker.distribution.manifest.v2+json"
	ha := internalHttp.Auth{HeaderKey: "Accept", HeaderValue: header, Method: "GET", Pass: dr.Pass, User: dr.User, URL: url}
	rc, err := ha.RequestAndResponseBody(nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()
	b, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	log.Trace(string(b))
	configDigest := gjson.Get(string(b), "config.digest").String()
	log.Debug(configDigest)
	digests := gjson.Get(string(b), "layers.#.digest").Array()

	s := make([]string, 0, len(digests))
	for _, digest := range digests {
		s = append(s, digest.String()+".tar")
	}

	log.Trace(digests)

	if err := dr.manifest(dir, configDigest, repo, tag, s); err != nil {
		return err
	}

	file := filepath.Join(dir, configDigest+".json")
	url = dr.Host + internalHttp.Version + repo + blobs + configDigest
	if err := dr.download(file, "application/vnd.docker.container.image.v1+json", url); err != nil {
		return err
	}
	return nil
}

func (dr *DockerRegistry) All() error {
	if err := dr.allLayers(); err != nil {
		return err
	}

	return nil
}

func (dr *DockerRegistry) allLayers() error {
	url := dr.Host + internalHttp.Version + uriCatalog
	repos, err := dr.json(url, "repositories")
	if err != nil {
		return err
	}
	log.Infof("repositories: '%s'", repos)

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		repoString := repo.String()
		go func(repoString string) {
			defer wg.Done()

			if err := dr.tags(repoString); err != nil {
				panic(err)
			}
		}(repoString)
	}
	wg.Wait()

	return nil
}

func (dr *DockerRegistry) Image(repo, tag string) error {
	dir := filepath.Join(dr.Dir, repo, tag)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	if err := dr.layers(dir, repo, tag); err != nil {
		return err
	}

	if err := dr.downloadManifest(dir, repo, tag); err != nil {
		return err
	}

	if err := dr.downloadConfig(dir, repo, tag); err != nil {
		return err
	}

	if err := tar(dir); err != nil {
		return err
	}

	return nil
}

func (dr *DockerRegistry) tags(repo string) error {
	url := dr.Host + internalHttp.Version + repo + uriTagsList
	tags, err := dr.json(url, "tags")
	if err != nil {
		return err
	}
	log.Infof("tags: '%s'", tags)

	var wg sync.WaitGroup
	for _, tag := range tags {
		wg.Add(1)
		tagString := tag.String()
		go func(tag string) {
			defer wg.Done()

			if err := dr.Image(repo, tag); err != nil {
				panic(err)
			}
		}(tagString)
	}
	wg.Wait()

	return nil
}

func (dr *DockerRegistry) layers(dir, repo, tag string) error {
	url := dr.Host + internalHttp.Version + repo + uriManifests + tag
	blobSums, err := dr.json(url, "fsLayers.#.blobSum")
	if err != nil {
		return err
	}
	log.Debugf("blobSums: '%s'", blobSums)

	var wg sync.WaitGroup

	for _, blobSum := range blobSums {
		wg.Add(1)
		blobSumString := blobSum.String()
		go func(blobSumString string) {
			defer wg.Done()

			if err := os.MkdirAll(dir, os.ModePerm); err != nil {
				panic(err)
			}

			file := filepath.Join(dir, blobSumString+".tar")
			if err := dr.downloadLayer(blobSumString, file, repo); err != nil {
				panic(err)
			}
		}(blobSumString)
	}
	wg.Wait()

	for _, blobSum := range blobSums {
		wg.Add(1)
		blobSumString := blobSum.String()
		go func(blobSumString string) {
			defer wg.Done()

			file := filepath.Join(dir, blobSumString+".tar")
			if err := validate(blobSumString, file); err != nil {
				panic(err)
			}
		}(blobSumString)
	}
	wg.Wait()

	return nil
}

func (dr *DockerRegistry) json(url, key string) ([]gjson.Result, error) {
	ha := internalHttp.Auth{HeaderKey: "Accept", Method: "GET", Pass: dr.Pass, User: dr.User, URL: url}
	rc, err := ha.RequestAndResponseBody(nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rc.Close(); err != nil {
			panic(err)
		}
	}()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	values := gjson.Get(string(b), key).Array()

	return values, err
}

func validate(expected, file string) error {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}

	checksum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	log.Debugf("%s vs. %s", checksum, expected)
	if checksum != expected {
		return fmt.Errorf("expected checksum: '%s', actual: '%s' (%s)", expected, checksum, file)
	}

	return nil
}

func (dr *DockerRegistry) manifest(dir, configDigest, repo, tag string, digests []string) error {
	re := regexp.MustCompile(`^http(s)?:\/\/(.*)$`)
	host := re.ReplaceAllString(dr.Host, `$2`)
	repoTag := host + "/" + repo + ":" + tag
	f := configDigest + ".json"

	b, err := json.Marshal([]manifestJSON{{Config: f, RepoTags: []string{repoTag}, Layers: digests}})
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "manifest.json"), b, 0o600); err != nil {
		return err
	}

	return nil
}

func tar(dir string) error {
	imageTar := "image.tar"
	m := make(map[string]string)

	if err := filepath.WalkDir(dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && d.Name() != imageTar {
				m[path] = ""
			}

			log.Trace(path)
			return nil
		}); err != nil {
		return err
	}

	files, err := archiverV4.FilesFromDisk(nil, m)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Clean(filepath.Join(dir, imageTar)))
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			panic(err)
		}
	}()
	format := archiverV4.CompressedArchive{
		Compression: archiverV4.Gz{},
		Archival:    archiverV4.Tar{},
	}
	err = format.Archive(context.Background(), out, files)
	if err != nil {
		return err
	}
	return nil
}
