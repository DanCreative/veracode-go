package veracode

import (
	"net/http"
	"time"

	"github.com/DanCreative/veracode-go/hmac"
	"golang.org/x/time/rate"
)

// veracodeTransport implements the http.RoundTripper interface and
// wraps either a provided http.RoundTripper or the http.DefaultTransport.
//
// veracodeTransport is responsible for client-side rate limiting as well as
// adding the Veracode HMAC Hash message to the Authorization header.
type veracodeTransport struct {
	r1        *rate.Limiter
	Key       string
	Secret    string
	Transport http.RoundTripper
}

// RoundTrip is required to implement the http.RoundTripper interface.
func (v *veracodeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the HMAC Hash message to the Authorization header
	bearer, err := hmac.CalculateAuthorizationHeader(req.URL, req.Method, v.Key, v.Secret)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", bearer)
	req.Header.Add("Content-Type", "application/json")

	// Wait for the limiter.
	err = v.r1.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	return v.transport().RoundTrip(req)
}

// transport checks if a custom http.RoundTripper was provided and returns it if it was and the http.DefaultTransport if it wasn't.
func (v *veracodeTransport) transport() http.RoundTripper {
	if v.Transport != nil {
		return v.Transport
	}
	return http.DefaultTransport
}

// Client() returns a new http.Client with this RateTransport object set as its Transport.
func (v *veracodeTransport) Client() *http.Client {
	return &http.Client{Transport: v}
}

// newTransport returns a new veracodeTransport.
func newTransport(rt http.RoundTripper, key, secret string, period time.Duration, burst int) *veracodeTransport {
	return &veracodeTransport{
		Transport: rt,
		r1:        rate.NewLimiter(rate.Every(period), burst),
		Key:       key,
		Secret:    secret,
	}
}
