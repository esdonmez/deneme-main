package pgo

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io"
	"net"
	"net/http"
	"time"
)

type ApiClient struct {
	Client             *http.Client
	DeadlineRetryCount int
}

type ApiResponse struct {
	Body       []byte
	StatusCode int
}

func NewHttpClient() *ApiClient {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   60 * time.Second,
			KeepAlive: 150 * time.Second,
		}).DialContext,
		MaxIdleConns:          300,
		MaxIdleConnsPerHost:   300,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	retryCount := 3
	retryClient := retryablehttp.NewClient()
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return 300 * time.Millisecond
	}
	retryClient.RetryMax = retryCount
	retryClient.HTTPClient.Transport = transport
	retryClient.HTTPClient.Timeout = 60 * time.Second
	retryClient.Logger = nil
	client := *retryClient.StandardClient()

	apiClient := ApiClient{
		Client:             &client,
		DeadlineRetryCount: retryCount,
	}

	return &apiClient
}

func (a *ApiClient) Get(url string) (ApiResponse, error) {
	req, _ := http.NewRequest(http.MethodGet, url, http.NoBody)
	return a.do(req)
}

func (a *ApiClient) do(req *http.Request) (ApiResponse, error) {
	response, err := a.Client.Do(req)

	if err != nil {
		return ApiResponse{},
			fmt.Errorf("an error occurred. Request: %s, err = %s", req.URL, err.Error())
	}
	defer response.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, response.Body)
	if err != nil {
		return ApiResponse{StatusCode: response.StatusCode},
			fmt.Errorf("an error occurred while reading body. Request: %s, err = %s", req.URL, err.Error())
	}

	return ApiResponse{
		StatusCode: response.StatusCode,
		Body:       buf.Bytes(),
	}, nil
}
