package downloader

import (
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	httpClient     *resty.Client
	httpClientOnce = false
)

func getHTTPClient() *resty.Client {
	if !httpClientOnce {
		httpClient = resty.New().
			SetTimeout(5 * time.Minute).
			SetRetryCount(3).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(30 * time.Second).
			SetTransport(&http.Transport{
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 5,
				IdleConnTimeout:     30 * time.Second,
			})
		httpClientOnce = true
	}
	return httpClient
}
