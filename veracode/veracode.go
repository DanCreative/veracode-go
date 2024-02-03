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
	urlEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	url := c.BaseURL.ResolveReference(urlEndpoint)

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
	default:
		fmt.Println("Body not matched")
	}
	return r
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error has occured during API Call: %d", resp.StatusCode)
	}
	return nil
}
