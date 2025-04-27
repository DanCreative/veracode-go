package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// CreateSandbox contains all of the fields required for creating and updating development sandboxes.
//
// Only the Name field is required.
type CreateSandbox struct {
	Name         string        `json:"name,omitempty"`
	AutoCreate   bool          `json:"auto_create,omitempty"` // If you are in the time-to-live mode, Automatically re-create the sandbox once the period expires. Documentation: https://docs.veracode.com/r/About_Sandbox_Data_Retention
	CustomFields []CustomField `json:"custom_fields,omitempty"`
}

type Sandbox struct {
	ApplicationGuid string        `json:"application_guid,omitempty"`
	Created         time.Time     `json:"created,omitempty"`
	CustomFields    []CustomField `json:"custom_fields,omitempty"`
	Guid            string        `json:"guid,omitempty"`
	Id              int           `json:"id,omitempty"`
	Modified        time.Time     `json:"modified,omitempty"`
	Name            string        `json:"name,omitempty"`
	OrganizationId  int           `json:"organization_id,omitempty"`
	OwnerUsername   string        `json:"owner_username,omitempty"`
}

type sandboxSearchResult struct {
	Embedded struct {
		Sandboxes []Sandbox `json:"sandboxes"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *sandboxSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *sandboxSearchResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListSandboxes takes an application GUID string and page options, and then returns a list of sandboxes for that application.
func (s *SandboxService) ListSandboxes(ctx context.Context, applicationGuid string, options PageOptions) ([]Sandbox, *Response, error) {
	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes", applicationGuid), http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var result sandboxSearchResult

	resp, err := s.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.Sandboxes, resp, nil
}

// GetSandbox takes an application GUID string and a sandbox GUID, and then returns the sandbox with the provided GUID.
func (s *SandboxService) GetSandbox(ctx context.Context, applicationGuid string, sandboxGuid string) (*Sandbox, *Response, error) {
	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes/%s", applicationGuid, sandboxGuid), http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var result Sandbox

	resp, err := s.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// CreateSandbox takes an application GUID and a CreateSandbox, and then creates a new sandbox for the provided application.
func (s *SandboxService) CreateSandbox(ctx context.Context, applicationGuid string, sandbox CreateSandbox) (*Sandbox, *Response, error) {
	byt, err := json.Marshal(sandbox)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes", applicationGuid), http.MethodPost, bytes.NewBuffer(byt))
	if err != nil {
		return nil, nil, err
	}

	var result Sandbox

	resp, err := s.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// UpdateSandbox takes an application GUID, a sandbox GUID and a CreateSandbox, and updates the existing sandbox with the new body.
func (s *SandboxService) UpdateSandbox(ctx context.Context, applicationGuid string, sandboxGuid string, sandbox CreateSandbox) (*Sandbox, *Response, error) {
	byt, err := json.Marshal(sandbox)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes/%s", applicationGuid, sandboxGuid), http.MethodPut, bytes.NewBuffer(byt))
	if err != nil {
		return nil, nil, err
	}

	var result Sandbox

	resp, err := s.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// DeleteSandbox takes an application GUID and a Sandbox GUID and deletes the sandbox with provide GUID.
func (s *SandboxService) DeleteSandbox(ctx context.Context, applicationGuid string, sandboxGuid string) (*Response, error) {
	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes/%s", applicationGuid, sandboxGuid), http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// PromoteSandbox promotes the latest scan in a sandbox, to a policy scan. setting deleteOnPromotion to true, will delete said scan once
// it has been promoted to the policy.
func (s *SandboxService) PromoteSandbox(ctx context.Context, applicationGuid string, sandboxGuid string, deleteOnPromotion bool) (*Sandbox, *Response, error) {
	req, err := s.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/applications/%s/sandboxes/%s/promote", applicationGuid, sandboxGuid), http.MethodPost, nil)
	if err != nil {
		return nil, nil, err
	}

	if deleteOnPromotion {
		req.URL.RawQuery = "delete_on_promote=true"
	}

	var result Sandbox

	resp, err := s.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}
