package veracode

import (
	"context"
	"net/http"
)

type ListCustomFieldOptions struct {
	Page int `url:"page"`
	Size int `url:"size,omitempty"`
}

type ApplicationCustomField struct {
	Name      string `json:"name,omitempty"`
	SortOrder int    `json:"sort_order,omitempty"`
}

type appCustomFieldSearchResult struct {
	Embedded struct {
		CustomFields []ApplicationCustomField `json:"app_custom_field_names"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *appCustomFieldSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *appCustomFieldSearchResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListCustomFields returns a list of the custom fields for the Application Profiles.
func (a *ApplicationService) ListCustomFields(ctx context.Context, options ListCustomFieldOptions) ([]ApplicationCustomField, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/custom_fields", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = QueryEncode(options)

	var results appCustomFieldSearchResult

	resp, err := a.Client.Do(req, &results)
	if err != nil {
		return nil, resp, err
	}

	return results.Embedded.CustomFields, resp, nil
}
