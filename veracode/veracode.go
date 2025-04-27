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
	"time"
)

type Client struct {
	baseURL    *url.URL
	rwMu       sync.RWMutex
	HttpClient *http.Client

	// Services used for talking to the different parts of the Veracode API
	common service

	Identity    *IdentityService    // See type for documentation.
	Application *ApplicationService // See type for documentation.
	Sandbox     *SandboxService     // See type for documentation.
	Healthcheck *HealthCheckService // See type for documentation.
}

type Response struct {
	*http.Response
	Page  PageMeta
	Links NavLinks
}

// Any struct that is used to unmarshal a collection of entities, needs to implement the CollectionResult interface in order for the page meta and navigational links
// to be set in the Response object.
type CollectionResult interface {
	GetLinks() NavLinks
	GetPageMeta() PageMeta
}

func NewClient(httpClient *http.Client, apiKey, apiSecret string) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// Wrap the transport provided in the http.Client with the veracodeTransport (which will handle rate limiting and authentication)
	httpClient.Transport = newTransport(httpClient.Transport, apiKey, apiSecret, time.Minute*1, 500)

	regionURL, err := GetRegionFromCredentials(apiKey)
	if err != nil {
		return nil, err
	}

	baseEndpoint, err := url.Parse(regionURL)
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
	c.Sandbox = (*SandboxService)(&c.common)
	c.Healthcheck = (*HealthCheckService)(&c.common)

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
	defer resp.Body.Close()

	err = checkStatus(resp)
	if err != nil {
		err = NewVeracodeError(resp)
		return newResponse(resp, nil), err
	}

	if body != nil {
		err = json.NewDecoder(resp.Body).Decode(body)
		if err != nil {
			return newResponse(resp, nil), err
		}
	}
	return newResponse(resp, body), nil
}

// UpdateCredentials is a method that allows the caller to update the credentials for the client
// after it has been initialized.
func (c *Client) UpdateCredentials(apiKey, apiSecret string) error {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	regionURL, err := GetRegionFromCredentials(apiKey)
	if err != nil {
		return err
	}

	v := c.HttpClient.Transport.(*veracodeTransport)
	v.Key, v.Secret = apiKey, apiSecret

	baseEndpoint, err := url.Parse(regionURL)
	if err != nil {
		return err
	}

	c.baseURL = baseEndpoint
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
		return fmt.Errorf("error has occurred during API Call: %d", resp.StatusCode)
	}
	return nil
}
