package downloader

import (
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	httpClient     *resty.Client
	httpClientOnce = false
)

func getHTTPClient() *resty.Client {
	if !httpClientOnce {
		client := resty.New().
			SetTimeout(5 * time.Minute).
			SetRetryCount(3).
			SetRetryWaitTime(5 * time.Second).
			SetRetryMaxWaitTime(30 * time.Second)

		proxyURL := os.Getenv("HTTPS_PROXY")
		if proxyURL == "" {
			proxyURL = os.Getenv("https_proxy")
		}
		if proxyURL == "" {
			proxyURL = os.Getenv("HTTP_PROXY")
		}
		if proxyURL == "" {
			proxyURL = os.Getenv("http_proxy")
		}

		if proxyURL != "" {
			if parsedURL, err := url.Parse(proxyURL); err == nil {
				client.SetTransport(&http.Transport{
					Proxy: http.ProxyURL(parsedURL),
				})
			}
		}

		httpClient = client
		httpClientOnce = true
	}
	return httpClient
}
