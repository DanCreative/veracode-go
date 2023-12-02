package veracode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	BaseURL    *url.URL
	HttpClient *http.Client

	// Services used for talking to the different parts of the SpaceTraders API
	common   service
	Identity *IdentityService
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

	c := &Client{
		BaseURL:    baseEndpoint,
		HttpClient: httpClient,
	}

	c.common.Client = c
	c.Identity = (*IdentityService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(ctx context.Context, endpoint string, method string, body io.Reader) (*http.Request, error) {
	url := c.BaseURL.JoinPath(endpoint)
	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, err
	}

	return req, err
}

func (c *Client) Do(req *http.Request, body any) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return resp, err
	}
	err = checkStatus(resp)
	if err != nil {
		return resp, err
	}

	if body != nil {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(body)
		if err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error has occured during API Call: %d", resp.StatusCode)
	}
	return nil
}
