package veracode

import (
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL *url.URL
	Client  *http.Client
}

func NewClient(region Region, httpClient *http.Client, apiKey, apiSecret string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	authTransport, err := NewAuthTransport(httpClient.Transport, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	httpClient.Transport = authTransport

	baseEndpoint, err := url.Parse(parseRegion(region))
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}

	return &Client{BaseURL: baseEndpoint, Client: httpClient}, nil
}
