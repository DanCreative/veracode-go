package veracode

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateTransport struct {
	r1        *rate.Limiter
	Transport http.RoundTripper
}

// RoundTrip is required to implement the http.RoundTripper interface.
// It limits the number of requests per period at the transport level.
func (r *RateTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	err := r.r1.Wait(req.Context())
	if err != nil {
		return nil, err
	}

	return r.transport().RoundTrip(req)
}

// transport checks if a custom http.RoundTripper was provided and returns it if it was and the http.DefaultTransport if it wasn't.
func (r *RateTransport) transport() http.RoundTripper {
	if r.Transport != nil {
		return r.Transport
	}
	return http.DefaultTransport
}

// Client() returns a new http.Client with this RateTransport object set as its Transport.
func (t *RateTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

// NewRateTransport takes a object that implements interface http.RoundTripper as well as the duration between calls and allowed burst and returns a new RateTransport.
func NewRateTransport(rt http.RoundTripper, period time.Duration, burst int) (*RateTransport, error) {
	r := &RateTransport{
		Transport: rt,
		r1:        rate.NewLimiter(rate.Every(period), burst),
	}

	return r, nil
}
