// converted https://raw.githubusercontent.com/moby/moby/master/contrib/download-frozen-image-v2.sh to golang

package pull

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	registry      = "https://registry-1.docker.io/v2/"
	imageManifest = `[
  {
    "Config": "{{ .ConfigDigest }}.json",
    "RepoTags": [
      "{{ .Repo }}:{{ .Tag }}"
    ],
    "Layers": [
	{{- range .DigestLayers}}
      "{{.}}/layer.tar",
	{{- end}}
    ]
  }
]
`
)

type Manifest struct {
	ConfigDigest, Repo, Tag string
	DigestLayers            []string
}

func DockerHub(dir, repo, tag string) error {
	m, err := manifest(repo, tag)
	if err != nil {
		return err
	}
	configDigest := gjson.Get(m, "config.digest").String()
	configSize := gjson.Get(m, "config.size").String()
	l := gjson.Get(m, "layers").Array()
	gjsonLayerDigests := gjson.Get(m, "layers.#.digest").Array()
	mediaType := gjson.Get(m, "mediaType").String()
	schemaVersion := gjson.Get(m, "schemaVersion").String()
	log.WithFields(log.Fields{
		"configDigest":  configDigest,
		"configSize":    configSize,
		"layers":        len(l),
		"layerDigests":  len(gjsonLayerDigests),
		"mediaType":     mediaType,
		"schemaVersion": schemaVersion,
	}).Debug("manifest")

	if err := config(configDigest, dir, repo, tag); err != nil {
		return err
	}

	var layerDigests []string
	for _, gjsonLayerDigest := range gjsonLayerDigests {
		digestWithoutSha256 := strings.TrimPrefix(gjsonLayerDigest.String(), "sha256:")
		layerDigests = append(layerDigests, digestWithoutSha256)
	}
	if err := layers(configDigest, dir, repo, tag, l, layerDigests); err != nil {
		return err
	}

	return nil
}

func layers(configDigest, dir, repo, tag string, layers []gjson.Result, layerDigests []string) error {
	dirRepoTag := filepath.Join(dir, repo, tag)
	f, err := os.Create(filepath.Join(dirRepoTag, "manifest.json"))
	if err != nil {
		return err
	}

	configDigestWithoutSha256 := strings.TrimPrefix(configDigest, "sha256:")
	t := template.Must(template.New("imageManifest").Parse(imageManifest))
	t.Execute(f, Manifest{ConfigDigest: configDigestWithoutSha256, Repo: repo, Tag: tag, DigestLayers: layerDigests})
	defer f.Close()

	fmt.Println(layers)
	layerID := ""
	for _, layer := range layers {
		digest := gjson.Get(layer.String(), "digest").String()
		mediaType := gjson.Get(layer.String(), "mediaType").String()
		size := gjson.Get(layer.String(), "size").String()
		log.WithFields(log.Fields{
			"digest":    digest,
			"mediaType": mediaType,
			"size":      size,
		}).Debug("layer")

		parentID := layerID
		s := parentID + "\n" + digest
		fmt.Println(s)
		h := sha256.New()
		h.Write([]byte(s))
		layerID = fmt.Sprintf("%x", h.Sum(nil))

		layerIdDir := filepath.Join(dirRepoTag, layerID)
		if err := os.MkdirAll(layerIdDir, 0o755); err != nil {
			return err
		}

		f, err := os.Create(filepath.Join(layerIdDir, "layer.tar"))
		if err != nil {
			return err
		}
		defer f.Close()

		req, err := http.NewRequest("GET", registry+repo+"/blobs/"+digest, nil)
		if err != nil {
			return err
		}
		t, err := token(repo)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+t)
		req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
		// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		sc := resp.StatusCode
		fmt.Println(sc)
		if sc == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			n, err := f.Write(bodyBytes)
			if err != nil {
				return err
			}
			fmt.Printf("wrote %d bytes\n", n)
			f.Sync()
		}

		// fmt.Println(gjson.Get(layer.String(), "mediaType"))
		// digest2 := gjson.Get(layer.String(), "digest").String()
		// fmt.Println(">>>>> digest2....: " + digest2)
		// // ddigestWithoutSha256 := strings.TrimPrefix(digest2, "sha256:")
		// // fmt.Println(">>>>> digestWithoutSha256:", ddigestWithoutSha256)

		fmt.Println("layerID:", layerID)
		if err := os.MkdirAll(layerIdDir, 0o755); err != nil {
			return err
		}
		f, err = os.Create(filepath.Join(layerIdDir, "VERSION"))
		if err != nil {
			return err
		}
		defer f.Close()
		n3, err := f.WriteString("1.0")
		if err != nil {
			return err
		}
		fmt.Printf("wrote %d bytes\n", n3)
		f.Sync()

		// 	//
		// 	// json
		// 	//
		// 	a := `
		// {
		//   "id": "` + layerID + `",
		//   "parent": "` + parentID + `",
		//   "created": "0001-01-01T00:00:00Z",
		//   "container_config": {
		//     "Hostname": "",
		//     "Domainname": "",
		//     "User": "",
		//     "AttachStdin": false,
		//     "AttachStdout": false,
		//     "AttachStderr": false,
		//     "Tty": false,
		//     "OpenStdin": false,
		//     "StdinOnce": false,
		//     "Env": null,
		//     "Cmd": null,
		//     "Image": "",
		//     "Volumes": null,
		//     "WorkingDir": "",
		//     "Entrypoint": null,
		//     "OnBuild": null,
		//     "Labels": null
		//   }
		// }`
		// 	f, err = os.Create(filepath.Join(layerIdDir, "json"))
		// 	if err != nil {
		// 		return err
		// 	}
		// 	defer f.Close()
		// 	n3, err = f.WriteString(a)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	fmt.Printf("wrote %d bytes\n", n3)
		// 	f.Sync()

	}
	return nil
}

func token(repo string) (string, error) {
	log.Debug("getting token...")

	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:" + repo + ":pull")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	token := ""
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString := string(bodyBytes)

		token = gjson.Get(bodyString, "token").String()
	}

	return token, nil
}

func config(configDigest, dir, repo, tag string) error {
	t, err := token(repo)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", registry+repo+"/blobs/"+configDigest, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
	// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		digestWithoutSha256 := strings.TrimPrefix(configDigest, "sha256:")

		dirRepoTag := filepath.Join(dir, repo, tag)
		if err := os.MkdirAll(dirRepoTag, 0o755); err != nil {
			return err
		}
		fp := filepath.Join(dirRepoTag, digestWithoutSha256+".json")
		f, err := os.Create(fp)
		if err != nil {
			return err
		}

		defer f.Close()

		n, err := f.Write(bodyBytes)
		if err != nil {
			return err
		}
		log.Debugf("number of bytes: '%d' written to filepath: '%s'", n, fp)
	}
	return nil
}

func manifest(repo, tag string) (string, error) {
	t, err := token(repo)
	if err != nil {
		return "", err
	}

	u := registry + repo + "/manifests/" + tag
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+t)
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
	// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyString := ""
	sc := resp.StatusCode
	if sc == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		bodyString = string(bodyBytes)

	} else {
		return "", fmt.Errorf("status code not 200, but: '%d'. URL: '%s'", sc, u)
	}

	if bodyString == "" {
		return "", fmt.Errorf("response should not be empty")
	}

	log.Trace(bodyString)
	return bodyString, nil
}
