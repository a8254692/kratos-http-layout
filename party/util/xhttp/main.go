package xhttp

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Client ...
type Client struct {
	*http.Client
}

// NewClient ...
func NewClient() *Client {
	var transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.DisableKeepAlives = true
	return &Client{
		Client: &http.Client{
			Transport: transport,
			Timeout:   10 * time.Second,
		},
	}
}

// HttpClient ...
func (c *Client) HttpClient() *http.Client {
	return c.Client
}

// SetDisableKeepAlive ...
func (c *Client) SetDisableKeepAlive(b bool) *Client {
	c.Transport.(*http.Transport).DisableKeepAlives = b
	return c
}

// SetTimeout ...
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.Timeout = timeout
	return c
}

// SetTLSClientConfig ...
func (c *Client) SetTLSClientConfig(tlsConfig *tls.Config) *Client {
	c.Transport.(*http.Transport).TLSClientConfig = tlsConfig
	return c
}
