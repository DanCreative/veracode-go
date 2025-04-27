package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ListBuOptions struct {
	SearchTerm  string `url:"search_term,omitempty"` // You can search for partial strings of the name.
	PageOptions        // can only sort by buName
}

type BusinessUnit struct {
	BuId       string  `json:"bu_id,omitempty"`
	BuLegacyId int     `json:"bu_legacy_id,omitempty"`
	BuName     string  `json:"bu_name,omitempty"`
	IsDefault  *bool   `json:"is_default,omitempty"`
	Teams      *[]Team `json:"teams,omitempty"`
}

// buSearchResult is required to decode the list of business units and search user response bodies.
type buSearchResult struct {
	Embedded struct {
		BusinessUnits []BusinessUnit `json:"business_units"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *buSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *buSearchResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListBusinessUnits returns a list of business units. A name can optionally be provided to search for BUs by name.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_list_bu
func (i *IdentityService) ListBusinessUnits(ctx context.Context, options ListBuOptions) ([]BusinessUnit, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/business_units", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var buResult buSearchResult

	resp, err := i.Client.Do(req, &buResult)
	if err != nil {
		return nil, resp, err
	}
	return buResult.Embedded.BusinessUnits, resp, err
}

// GetBusinessUnit returns the BusinessUnit with the provided buId.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_bu_info
func (i *IdentityService) GetBusinessUnit(ctx context.Context, buId string) (*BusinessUnit, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/business_units/"+buId, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var getBu BusinessUnit

	resp, err := i.Client.Do(req, &getBu)
	if err != nil {
		return nil, resp, err
	}
	return &getBu, resp, err
}

// CreateBusinessUnit creates a new bu using the provided BusinessUnit object.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_create_bu
func (i *IdentityService) CreateBusinessUnit(ctx context.Context, bu *BusinessUnit) (*BusinessUnit, *Response, error) {
	buf, err := json.Marshal(bu)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/business_units", http.MethodPost, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var newBu BusinessUnit
	resp, err := i.Client.Do(req, &newBu)
	if err != nil {
		return nil, resp, err
	}

	return &newBu, resp, nil
}

// UpdateBusinessUnit updates a specific bu and sets nulls to fields not in the request (if the database allows it) unless partial is set to true.
// If incremental is set to true, any values in the teams list will be added to the bu's teams instead of replacing them.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_update_bu
//   - https://docs.veracode.com/r/c_identity_add_team_bu
func (i *IdentityService) UpdateBusinessUnit(ctx context.Context, bu *BusinessUnit, options UpdateOptions) (*BusinessUnit, *Response, error) {
	buf, err := json.Marshal(bu)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/business_units/"+bu.BuId, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var updatedBu BusinessUnit
	resp, err := i.Client.Do(req, &updatedBu)
	if err != nil {
		return nil, resp, err
	}

	return &updatedBu, resp, nil
}

// DeleteBusinessUnit deletes a bu from the Veracode API using the provided buId.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_delete_bu
func (i *IdentityService) DeleteBusinessUnit(ctx context.Context, buId string) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/business_units/"+buId, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
