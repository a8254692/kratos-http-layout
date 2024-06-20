package util

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"net/http"
	"syscall"
	"time"
)

var (
	// DefaultHttpClient ...
	DefaultHttpClient *http.Client
	// ShortConnHttpClient ...
	ShortConnHttpClient *http.Client
	// DefaultRestyClient ...
	DefaultRestyClient *resty.Client
)

func init() {
	var transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxConnsPerHost = 100
	transport.MaxIdleConnsPerHost = 100

	// 默认全局 http client
	DefaultHttpClient = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// 全局短连接 http client
	tr := transport.Clone()
	tr.DisableKeepAlives = true
	ShortConnHttpClient = &http.Client{
		Transport: tr,
		Timeout:   10 * time.Second,
	}

	// 默认 RestyClient
	DefaultRestyClient = resty.New().
		SetTimeout(time.Second * 10).
		SetTransport(transport).
		SetRetryCount(5).
		SetRetryWaitTime(200 * time.Millisecond).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			// NOTE: Too Many Requests || connection reset by peer || broken pipe
			return response.StatusCode() == http.StatusTooManyRequests ||
				errors.Is(err, syscall.ECONNRESET) ||
				errors.Is(err, syscall.EPIPE)
		})
}
