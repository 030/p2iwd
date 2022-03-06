package push

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	internalHttp "github.com/030/p2iwd/internal/pkg/http"
	log "github.com/sirupsen/logrus"
)

type DockerImage struct {
	name, tag string
}

type DockerRegistry struct {
	Dir, Host, Pass, User string
}

func (dr *DockerRegistry) location(dockerImageName string) (string, error) {
	url := dr.Host + internalHttp.Version + dockerImageName + "/blobs/uploads/"
	log.Debugf("URL: '%s'", url)
	ha := internalHttp.Auth{Method: "POST", Pass: dr.Pass, User: dr.User, URL: url}
	resp, err := ha.RequestAndResponse(nil)
	if err != nil {
		return "", err
	}

	location := resp.Header.Get("Location")
	log.Debug(location)
	return location, nil
}

func checksum(file *os.File) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	checksum := fmt.Sprintf("sha256:%x", h.Sum(nil))
	return checksum, nil
}

func (dr *DockerRegistry) dockerImageNameAndTag(path string) (DockerImage, error) {
	res := strings.ReplaceAll(path, dr.Dir, "")
	re := regexp.MustCompile(`^(.*)/(v?[0-9].*)/.*$`)
	match := re.FindStringSubmatch(res)
	dockerImageName := match[1]
	log.Debugf("dockerImageName: '%s'", dockerImageName)
	tag := match[2]
	length := len(match)
	log.Debugf("tag: '%s'. Len: '%d'", tag, length)
	if length != 3 {
		return DockerImage{}, fmt.Errorf("is not three: '%d'", length)
	}
	return DockerImage{name: dockerImageName, tag: tag}, nil
}

func (dr *DockerRegistry) All() error {
	if err := filepath.WalkDir(dr.Dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && d.Name() != "manifest.json" && (filepath.Ext(path) == ".json" || filepath.Ext(path) == ".tar") {
				log.Debugf("found path: '%s'", path)
				f, err := os.Open(filepath.Clean(path))
				if err != nil {
					return err
				}
				defer func() {
					if err := f.Close(); err != nil {
						panic(err)
					}
				}()

				log.Debugf("calculating checksum for path: '%s'", path)
				cs, err := checksum(f)
				if err != nil {
					return err
				}

				log.Debugf("determine dockerImageName and tag for path: '%s'", path)
				dinat, err := dr.dockerImageNameAndTag(path)
				if err != nil {
					return err
				}

				l, err := dr.location(dinat.name)
				if err != nil {
					return err
				}
				u := url.URL{Host: l}
				log.Debugf("check whether location: '%s' is a absolute URL. Outcome: '%t'", l, u.IsAbs())
				uploadURL := l + "?digest=" + cs
				if !u.IsAbs() {
					uploadURL = dr.Host + uploadURL
				}
				log.Debugf("URL: '%s'", uploadURL)
				log.Debugf("Filename: '%s'", d.Name())

				headerValue := "application/vnd.docker.image.rootfs.diff.tar.gzip"
				if d.Name() == "upload-manifest.json" {
					headerValue = "application/vnd.docker.distribution.manifest.v2+json"
					uploadURL = dr.Host + internalHttp.Version + dinat.name + "/manifests/" + dinat.tag
					log.Debugf("Output regex extract:'%s' '%s' -> '%s'", dinat.name, dinat.tag, uploadURL)
				}

				ha := internalHttp.Auth{HeaderKey: "Content-Type", HeaderValue: headerValue, Method: "PUT", Pass: dr.Pass, User: dr.User, URL: uploadURL}
				rc, err := ha.RequestAndResponseBody(f)
				if err != nil {
					return err
				}
				defer rc.Close()
				b, err := io.ReadAll(rc)
				if err != nil {
					return err
				}
				log.Debug(string(b))
			}
			return nil
		}); err != nil {
		return err
	}

	return nil
}
