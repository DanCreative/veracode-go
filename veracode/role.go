package veracode

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

// RoleUser struct contains the fields that are return as part of the user aggregate.
type RoleUser struct {
	RoleDescription string `json:"role_description,omitempty"`
	RoleId          string `json:"role_id,omitempty"`
	RoleName        string `json:"role_name,omitempty"`
}

type Role struct {
	IsApi               bool   `json:"is_api,omitempty"`
	IsScanType          bool   `json:"is_scan_type,omitempty"`
	TeamAdminManageable bool   `json:"team_admin_manageable,omitempty"`
	RoleDescription     string `json:"role_description,omitempty"`
	RoleId              string `json:"role_id,omitempty"`
	RoleName            string `json:"role_name,omitempty"`
	RoleLegacyId        int    `json:"role_legacy_id,omitempty"`
	// BACKLOG: Add remaining fields for model as required.
}

// roleSearchResult is required to decode the list of roles and search roles response bodies.
type roleSearchResult struct {
	Embedded struct {
		Roles []Role `json:"roles"`
	} `json:"_embedded"`
	Links navLinks `json:"_links"`
	Page  pageMeta `json:"page"`
}

// ListRoles takes a PageOptions and returns a list of roles.
//
// Veracode API documentation: https://docs.veracode.com/r/Listing_All_Roles_in_an_Organization_with_the_Identity_API.
func (i *IdentityService) ListRoles(ctx context.Context, options PageOptions) ([]Role, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/roles", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = values.Encode()

	var rolesResult roleSearchResult

	resp, err := i.Client.Do(req, &rolesResult)
	if err != nil {
		return nil, resp, err
	}
	return rolesResult.Embedded.Roles, resp, err
}
