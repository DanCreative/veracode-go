package veracode

import (
	"net/http"

	vmac "github.com/DanCreative/veracode-hmac-go"
)

type AuthTransport struct {
	key       string
	secret    string
	Transport http.RoundTripper
}

// RoundTrip is required to implement the http.RoundTripper interface.
// It gets the Authorization value using the github.com/DanCreative/veracode-hmac-go module.
func (t *AuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	bearer, err := vmac.CalculateAuthorizationHeader(req.URL, req.Method, t.key, t.secret)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	return t.transport().RoundTrip(req)
}

// transport checks if a custom http.RoundTripper was provided and returns it if it was and the http.DefaultTransport if it wasn't.
func (t *AuthTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

// Client() returns a new http.Client with this AuthTransport object set as its Transport.
func (t *AuthTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// NewAuthTransport takes a object that implements interface http.RoundTripper as well as a Veracode API Key and Secret and returns a new AuthTransport.
func NewAuthTransport(rt http.RoundTripper, key, secret string) (AuthTransport, error) {
	t := AuthTransport{
		Transport: rt,
		key:       key,
		secret:    secret,
	}

	return t, nil
}
