package http

import (
	"encoding/json"
	"fmt"
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

type ErrorResponse struct {
	Errors []struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Detail  string `json:"detail"`
	} `json:"errors"`
}

func (a *Auth) RequestAndResponse(body io.Reader) (*http.Response, error) {
	log.Tracef("trying to authenticate to host: '%s'...", a.URL)
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
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)

		var response ErrorResponse
		if err := json.Unmarshal([]byte(bodyString), &response); err != nil {
			fmt.Println("Error:", err)
			return nil, fmt.Errorf("error")
		}
		errorMessage := ""
		if len(response.Errors) > 0 {
			errorMessage = response.Errors[0].Detail
			fmt.Println("Message:", errorMessage)
		} else {
			fmt.Println("No errors found in the response")
		}

		return nil, fmt.Errorf("statuscode was not 200, but: '%d' with response: %v", resp.StatusCode, errorMessage)
	}

	return resp, nil
}

func (a *Auth) RequestAndResponseBody(body io.Reader) (io.ReadCloser, error) {
	log.Tracef(">>>>>>>>>>>>>>>>CP3<<<<<<<<<<<<<<<<<<<<<")
	resp, err := a.RequestAndResponse(body)
	log.Tracef(">>>>>>>>>>>>>>>>CP3a<<<<<<<<<<<<<<<<<<<<<")
	if err != nil {
		return nil, err
	}
	log.Tracef(">>>>>>>>>>>>>>>>CP3b<<<<<<<<<<<<<<<<<<<<<")
	return resp.Body, nil
}
