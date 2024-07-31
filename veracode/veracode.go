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
	common service

	// You can use the Identity Service to manage the administrative configuration for your organization that is in the Veracode Platform.
	// For more information: https://docs.veracode.com/r/c_identity_intro.
	//
	// Currently supports V2 of the Identity API
	Identity *IdentityService
	// You can use the Applications API to quickly access information about your Veracode applications.
	// For more information, review the documentation: https://docs.veracode.com/r/c_apps_intro
	//
	// Currently supports V1 of the Applications API
	Application *ApplicationService
}

type Response struct {
	*http.Response
	Page  pageMeta
	Links navLinks
}

// Any struct that is used to unmarshal a collection of entities, needs to implement the CollectionResult interface in order for the page meta and navigational links
// to be set in the Response object.
type CollectionResult interface {
	GetLinks() navLinks
	GetPageMeta() pageMeta
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
	c.Application = (*ApplicationService)(&c.common)

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

// Do is a helper method that executes the provided http.Request and marshals the JSON response body
// into either the provided any object or into an error if an error occurred.
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

	if body == nil {
		return r
	}

	// if the body object provided implements the CollectionResult interface, meaning that it is a collection of entities, set the Links and Page Meta on the
	// Response object.
	if value, ok := body.(CollectionResult); ok {
		r.Links = value.GetLinks()
		r.Page = value.GetPageMeta()
	}

	return r
}

func checkStatus(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error has occured during API Call: %d", resp.StatusCode)
	}
	return nil
}
