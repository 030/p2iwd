package http

import (
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
	log "github.com/sirupsen/logrus"
)

const (
	Version = "/v2/"
)

type Auth struct {
	HeaderKey, HeaderValue, Method, Pass, URL, User string
}

func (a *Auth) RequestAndResponse(body io.Reader, token string) (*http.Response, error) {
	log.Info(">>>>>>>>>>>>>>>>>>>>>>>>>>>", token)
	req, err := http.NewRequest(a.Method, a.URL, body)
	if err != nil {
		return nil, err
	}
	if a.HeaderValue != "" {
		req.Header.Set(a.HeaderKey, a.HeaderValue)
		// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.list.v2+json")
		// req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v1+json")
		req.Header.Set("Authorization", "Bearer "+token)
	}
	// req.SetBasicAuth(a.User, a.Pass)

	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient()
	resp, err := standardClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Debug(resp.Status)
	log.Infof("url: '%s'. StatusCode: '%d'", a.URL, resp.StatusCode)

	return resp, nil
}

func (a *Auth) RequestAndResponseBody(body io.Reader, token string) (io.ReadCloser, error) {
	resp, err := a.RequestAndResponse(body, token)
	if err != nil {
		return nil, err
	}
	log.Trace(resp.StatusCode)

	return resp.Body, nil
}
