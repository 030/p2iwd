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

func (a *Auth) RequestAndResponse(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(a.Method, a.URL, body)
	if err != nil {
		return nil, err
	}
	if a.HeaderValue != "" {
		req.Header.Set(a.HeaderKey, a.HeaderValue)
	}
	req.SetBasicAuth(a.User, a.Pass)

	retryClient := retryablehttp.NewClient()
	retryClient.Logger = nil
	retryClient.RetryMax = 5
	standardClient := retryClient.StandardClient()
	resp, err := standardClient.Do(req)
	if err != nil {
		return nil, err
	}
	log.Debug(resp.Status)
	log.Debug(resp.StatusCode)

	return resp, nil
}

func (a *Auth) RequestAndResponseBody(body io.Reader) (io.ReadCloser, error) {
	resp, err := a.RequestAndResponse(body)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
