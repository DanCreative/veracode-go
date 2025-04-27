package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ListCollectionOptions struct {
	Page         int    `url:"page,omitempty"`
	Size         int    `url:"size,omitempty"`
	Name         string `url:"name,omitempty"`          // Filter collections by name (partial match)
	BusinessUnit string `url:"business_unit,omitempty"` // Filter collections by business unit name (partial match)
	Tag          string `url:"tag,omitempty"`           // Filter by tags
	// CustomFieldNames and CustomFieldValues need to both be set together.
	// You can use the AddCustomFieldOption method to set/update these fields.
	CustomFieldNames  []string `url:"custom_field_names,omitempty"`
	CustomFieldValues []string `url:"custom_field_values,omitempty"`
}

// AddCustomFieldOption sets the customFieldName and customFieldValue attributes on the ListApplicationOptions.
// To identify application profiles with any value for a specific custom field, enter the URL-encoded wildcard value %25 for customFieldValue.
//
// Documentation Reference: https://docs.veracode.com/r/r_applications_custom_field
func (l *ListCollectionOptions) AddCustomFieldOption(customFieldName, customFieldValue string) {
	if len(l.CustomFieldNames) < 1 || l.CustomFieldNames == nil {
		l.CustomFieldNames = []string{customFieldName}
	} else {
		l.CustomFieldNames = append(l.CustomFieldNames, customFieldName)
	}

	if len(l.CustomFieldValues) < 1 || l.CustomFieldValues == nil {
		l.CustomFieldValues = []string{customFieldValue}
	} else {
		l.CustomFieldValues = append(l.CustomFieldValues, customFieldName)
	}
}

type Collection struct {
	Assets       []CollectionAsset        `json:"asset_infos,omitempty"`
	BusinessUnit *ApplicationBusinessUnit `json:"business_unit,omitempty"`
	CustomFields []CustomField            `json:"custom_fields,omitempty"`
	Description  string                   `json:"description,omitempty"`
	Name         string                   `json:"name,omitempty"`
	Guid         string                   `json:"guid,omitempty"`
	Restricted   *bool                    `json:"restricted,omitempty"`
}

type CollectionAsset struct {
	Type string `json:"type,omitempty"`
	Guid string `json:"guid,omitempty"`
}

type collectionSearchResult struct {
	Embedded struct {
		Collections []Collection `json:"collections"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *collectionSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *collectionSearchResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListCollections returns []Collection using provided CollectionListOptions.
func (c *ApplicationService) ListCollections(ctx context.Context, options ListCollectionOptions) ([]Collection, *Response, error) {
	req, err := c.Client.NewRequest(ctx, "/appsec/v1/collections", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var results collectionSearchResult

	resp, err := c.Client.Do(req, &results)
	if err != nil {
		return nil, resp, err
	}

	return results.Embedded.Collections, resp, nil
}

// CreateCollection creates a new collection using the provided Collection.
func (c *ApplicationService) CreateCollection(ctx context.Context, collection Collection) (*Collection, *Response, error) {
	byt, err := json.Marshal(&collection)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.Client.NewRequest(ctx, "/appsec/v1/collections", http.MethodPost, bytes.NewBuffer(byt))
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Client.Do(req, &collection)
	if err != nil {
		return nil, resp, err
	}

	return &collection, resp, err
}

// UpdateCollection updates a collection with collectionId using provided collection.
func (c *ApplicationService) UpdateCollection(ctx context.Context, collection Collection) (*Collection, *Response, error) {
	byt, err := json.Marshal(&collection)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.Client.NewRequest(ctx, "/appsec/v1/collections/"+collection.Guid, http.MethodPut, bytes.NewBuffer(byt))
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.Client.Do(req, &collection)
	if err != nil {
		return nil, resp, err
	}

	return &collection, resp, nil
}

// GetCollection retrieves a collection with the provided collectionGuid.
func (a *ApplicationService) GetCollection(ctx context.Context, collectionGuid string) (*Collection, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/collections/"+collectionGuid, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var collection Collection

	resp, err := a.Client.Do(req, &collection)
	if err != nil {
		return nil, resp, err
	}

	return &collection, resp, nil
}

// GetCollection deletes a collection with the provided collectionGuid.
func (a *ApplicationService) DeleteCollection(ctx context.Context, collectionGuid string) (*Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/collections/"+collectionGuid, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}
	resp, err := a.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
