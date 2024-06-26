package veracode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Client struct {
	baseURL    *url.URL
	rwMu       sync.RWMutex
	HttpClient *http.Client

	// Services used for talking to the different parts of the Veracode API
	common   service
	Identity *IdentityService
}

type Response struct {
	*http.Response
	Page  pageMeta
	Links navLinks
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

	baseEndpoint, err := url.Parse(string(region))
	if err != nil {
		return nil, err
	}

	if !strings.HasSuffix(baseEndpoint.Path, "/") {
		baseEndpoint.Path += "/"
	}

	c := &Client{
		baseURL:    baseEndpoint,
		HttpClient: httpClient,
	}

	c.common.Client = c
	c.Identity = (*IdentityService)(&c.common)

	return c, nil
}

func (c *Client) NewRequest(ctx context.Context, endpoint string, method string, body io.Reader) (*http.Request, error) {
	urlEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	url := c.baseURL.ResolveReference(urlEndpoint)

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, err
	}

	return req, err
}

func (c *Client) Do(req *http.Request, body any) (*Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return newResponse(resp, nil), err
	}
	err = checkStatus(resp)
	if err != nil {
		err = NewVeracodeError(resp)
		return newResponse(resp, nil), err
	}

	if body != nil {
		defer resp.Body.Close()
		err = json.NewDecoder(resp.Body).Decode(body)
		if err != nil {
			return newResponse(resp, nil), err
		}
	}
	return newResponse(resp, body), nil
}

// SetBaseURL takes a Region and updates the Client's baseURL. It write locks the sync.RWMutex to prevent any concurrent reads while writing (just in case).
func (c *Client) SetBaseURL(region Region) error {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	new, err := url.Parse(string(region))
	if err != nil {
		return fmt.Errorf("error occured while trying to parse new base URL: %s", region)
	}

	c.baseURL = new
	return nil
}

func newResponse(response *http.Response, body any) *Response {
	r := &Response{
		Response: response,
	}

	switch value := body.(type) {
	case *userSearchResult:
		r.Links = value.Links
		r.Page = value.Page
	case *teamSearchResult:
		r.Links = value.Links
		r.Page = value.Page
	case *roleSearchResult:
		r.Links = value.Links
		r.Page = value.Page
	case *buSearchResult:
		r.Links = value.Links
		r.Page = value.Page
	}
	return r
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error has occured during API Call: %d", resp.StatusCode)
	}
	return nil
}
