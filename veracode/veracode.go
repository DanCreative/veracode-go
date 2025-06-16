package veracode

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type Client struct {
	baseRestURL *url.URL
	baseXmlURL  *url.URL
	rwMu        sync.RWMutex
	HttpClient  *http.Client

	// Services used for talking to the different parts of the Veracode API
	common service

	Identity    *IdentityService    // See type for documentation.
	Application *ApplicationService // See type for documentation.
	Sandbox     *SandboxService     // See type for documentation.
	Healthcheck *HealthCheckService // See type for documentation.
	UploadXML   *UploadXMLService   // See type for documentation.
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

	region, err := GetRegionFromCredentials(apiKey)
	if err != nil {
		return nil, err
	}

	c := &Client{
		HttpClient: httpClient,
	}

	setBaseURLs(c, region)

	c.common.Client = c
	c.Identity = (*IdentityService)(&c.common)
	c.Application = (*ApplicationService)(&c.common)
	c.Sandbox = (*SandboxService)(&c.common)
	c.Healthcheck = (*HealthCheckService)(&c.common)
	c.UploadXML = (*UploadXMLService)(&c.common)

	return c, nil
}

func setBaseURLs(c *Client, r Region) {
	for _, apiType := range []string{"rest", "xml"} {
		baseEndpoint, _ := url.Parse(r[apiType])

		if !strings.HasSuffix(baseEndpoint.Path, "/") {
			baseEndpoint.Path += "/"
		}

		switch apiType {
		case "rest":
			c.baseRestURL = baseEndpoint
		case "xml":
			c.baseXmlURL = baseEndpoint
		}
	}
}

// NewRequest is a helper method that creates a new request using the [Client]'s settings.
//
// By default, NewRequest will set the base URL to the REST variant, the caller can optionally provide shouldUseXML
// to switch to the XML base URL.
func (c *Client) NewRequest(ctx context.Context, endpoint string, method string, body io.Reader, shouldUseXML ...bool) (*http.Request, error) {
	urlEndpoint, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	c.rwMu.RLock()
	defer c.rwMu.RUnlock()

	var url *url.URL

	if len(shouldUseXML) > 0 && shouldUseXML[0] {
		url = c.baseXmlURL.ResolveReference(urlEndpoint)
	} else {
		url = c.baseRestURL.ResolveReference(urlEndpoint)
	}

	req, err := http.NewRequestWithContext(ctx, method, url.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

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

	contentType := resp.Header.Get("Content-Type")

	if body != nil {
		switch contentType {
		case "application/json", "application/json;charset=UTF-8":
			return doJson(resp, body)

		case "text/xml":
			return doXml(resp, body)

		default:
			return newResponse(resp, nil), fmt.Errorf("response header: 'Content-Type' contains unsupported value: %s", contentType)
		}
	}

	return newResponse(resp, body), nil
}

// doJson handles responses from the JSON APIs.
//
// Note: Only doJson contains the [checkStatus] function. That is because the Veracode XML APIs return 200 codes
// even if there is an error. Therefore the function is only applicable here and not on [doXml]
func doJson(resp *http.Response, body any) (*Response, error) {
	err := checkStatus(resp)
	if err != nil {
		err = NewVeracodeError(resp)
		return newResponse(resp, nil), err
	}
	err = json.NewDecoder(resp.Body).Decode(body)
	if err != nil {
		return newResponse(resp, nil), err
	}

	return newResponse(resp, body), nil
}

// doXml handles responses from the XML APIs.
//
// Note: Only doJson contains the [checkStatus] function. That is because the Veracode XML APIs return 200 codes
// even if there is an error. Therefore this function does not check the status code, it checks the first token
// in the XML document to determine whether an error occurred.
func doXml(resp *http.Response, body any) (*Response, error) {
	decoder := xml.NewDecoder(resp.Body)

	// Decode to the first StartElement and determine if its name equals "error" or not.
	// If its name equals "error", unmarshal the Response.Body into the Veracode Error and return that.
	// Otherwise, unmarshal the Response.Body into the provided body.
	for {
		t, err := decoder.Token()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return newResponse(resp, nil), err
		}

		if token, ok := t.(xml.StartElement); ok {
			if token.Name.Local == "error" {
				verr := Error{Code: resp.StatusCode, Endpoint: resp.Request.URL.Path}

				err = decoder.DecodeElement(&verr, &token)
				if err != nil {
					return newResponse(resp, nil), err
				}

				return newResponse(resp, nil), verr
			} else {
				err = decoder.DecodeElement(body, &token)
				if err != nil {
					return newResponse(resp, nil), err
				}

				return newResponse(resp, body), nil
			}
		}
	}

	return newResponse(resp, nil), nil
}

// UpdateCredentials is a method that allows the caller to update the credentials for the client
// after it has been initialized.
func (c *Client) UpdateCredentials(apiKey, apiSecret string) error {
	c.rwMu.Lock()
	defer c.rwMu.Unlock()

	region, err := GetRegionFromCredentials(apiKey)
	if err != nil {
		return err
	}

	v := c.HttpClient.Transport.(*veracodeTransport)
	v.Key, v.Secret = apiKey, apiSecret

	setBaseURLs(c, region)

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
