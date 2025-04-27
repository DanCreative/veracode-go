package veracode

import (
	"context"
	"net/http"
)

// GetStatus is a lightweight check that indicates whether the authentication services are operational.
//
// If GetStatus does not return an error, then everything is operational.
//
// Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-healthcheck_api_specification/1.0#/Healthcheck%20APIs/get_healthcheck_status
func (h *HealthCheckService) GetStatus(ctx context.Context) (*Response, error) {
	req, err := h.Client.NewRequest(ctx, "/healthcheck/status", http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
