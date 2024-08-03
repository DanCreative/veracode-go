package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type BusinessCriticality string
type ScanType string
type PolicyCompliance string
type ScanStatus string

const (
	VERY_HIGH BusinessCriticality = "VERY_HIGH"
	HIGH      BusinessCriticality = "HIGH"
	MEDIUM    BusinessCriticality = "MEDIUM"
	LOW       BusinessCriticality = "LOW"
	VERY_LOW  BusinessCriticality = "VERY_LOW"

	STATIC  ScanType = "STATIC"
	DYNAMIC ScanType = "DYNAMIC"
	MANUAL  ScanType = "MANUAL"

	PASSED           PolicyCompliance = "PASSED"
	CONDITIONAL_PASS PolicyCompliance = "CONDITIONAL_PASS"
	DID_NOT_PASS     PolicyCompliance = "DID_NOT_PASS"
	NOT_ASSESSED     PolicyCompliance = "NOT_ASSESSED"
	VENDOR_REVIEW    PolicyCompliance = "VENDOR_REVIEW"
	DETERMINING      PolicyCompliance = "DETERMINING"

	CREATED                                ScanStatus = "CREATED"
	UNPUBLISHED                            ScanStatus = "UNPUBLISHED"
	DELETED                                ScanStatus = "DELETED"
	PARTIAL_PUBLISH                        ScanStatus = "PARTIAL_PUBLISH"
	PARTIAL_UNPUBLISH                      ScanStatus = "PARTIAL_UNPUBLISH"
	INCOMPLETE                             ScanStatus = "INCOMPLETE"
	SCAN_SUBMITTED                         ScanStatus = "SCAN_SUBMITTED"
	IN_QUEUE                               ScanStatus = "IN_QUEUE"
	STOPPING                               ScanStatus = "STOPPING"
	PAUSING                                ScanStatus = "PAUSING"
	IN_PROGRESS                            ScanStatus = "IN_PROGRESS"
	ANALYSIS_ERRORS                        ScanStatus = "ANALYSIS_ERRORS"
	SCAN_CANCELED                          ScanStatus = "SCAN_CANCELED"
	INTERNAL_REVIEW                        ScanStatus = "INTERNAL_REVIEW"
	VERIFYING_RESULTS                      ScanStatus = "VERIFYING_RESULTS"
	SUBMITTED_FOR_NTO_PRE_SCAN             ScanStatus = "SUBMITTED_FOR_NTO_PRE_SCAN"
	SUBMITTED_FOR_DYNAMIC_PRE_SCAN         ScanStatus = "SUBMITTED_FOR_DYNAMIC_PRE_SCAN"
	PRE_SCAN_FAILED                        ScanStatus = "PRE_SCAN_FAILED"
	READY_TO_SUBMIT                        ScanStatus = "READY_TO_SUBMIT"
	NTO_PENDING_SUBMISSION                 ScanStatus = "NTO_PENDING_SUBMISSION"
	PRE_SCAN_COMPLETE                      ScanStatus = "PRE_SCAN_COMPLETE"
	MODULE_SELECTION_REQUIRED              ScanStatus = "MODULE_SELECTION_REQUIRED"
	PENDING_VENDOR_ACCEPTANCE              ScanStatus = "PENDING_VENDOR_ACCEPTANCE"
	SHOW_OSRDB                             ScanStatus = "SHOW_OSRDB"
	PUBLISHED                              ScanStatus = "PUBLISHED"
	PUBLISHED_TO_VENDOR                    ScanStatus = "PUBLISHED_TO_VENDOR"
	PUBLISHED_TO_ENTERPRISE                ScanStatus = "PUBLISHED_TO_ENTERPRISE"
	PENDING_ACCOUNT_APPROVAL               ScanStatus = "PENDING_ACCOUNT_APPROVAL"
	PENDING_LEGAL_AGREEMENT                ScanStatus = "PENDING_LEGAL_AGREEMENT"
	SCAN_IN_PROGRESS                       ScanStatus = "SCAN_IN_PROGRESS"
	SCAN_IN_PROGRESS_PARTIAL_RESULTS_READY ScanStatus = "SCAN_IN_PROGRESS_PARTIAL_RESULTS_READY"
	PROMOTE_IN_PROGRESS                    ScanStatus = "PROMOTE_IN_PROGRESS"
	PRE_SCAN_CANCELED                      ScanStatus = "PRE_SCAN_CANCELED"
	NTO_PRE_SCAN_CANCELED                  ScanStatus = "NTO_PRE_SCAN_CANCELED"
	SCAN_HELD_APPROVAL                     ScanStatus = "SCAN_HELD_APPROVAL"
	SCAN_HELD_LOGIN_INSTRUCTIONS           ScanStatus = "SCAN_HELD_LOGIN_INSTRUCTIONS"
	SCAN_HELD_LOGIN                        ScanStatus = "SCAN_HELD_LOGIN"
	SCAN_HELD_INSTRUCTIONS                 ScanStatus = "SCAN_HELD_INSTRUCTIONS"
	SCAN_HELD_HOLDS_FINISHED               ScanStatus = "SCAN_HELD_HOLDS_FINISHED"
	SCAN_REQUESTED                         ScanStatus = "SCAN_REQUESTED"
	TIMEFRAMEPENDING_ID                    ScanStatus = "TIMEFRAMEPENDING_ID"
	PAUSED_ID                              ScanStatus = "PAUSED_ID"
	STATIC_VALIDATING_UPLOAD               ScanStatus = "STATIC_VALIDATING_UPLOAD"
	PUBLISHED_TO_ENTERPRISEINT             ScanStatus = "PUBLISHED_TO_ENTERPRISEINT"
)

type Application struct {
	Guid    string             `json:"guid,omitempty"`
	Profile ApplicationProfile `json:"profile,omitempty"`
}

type ApplicationProfile struct {
	Name           string                   `json:"name,omitempty"`
	Tags           string                   `json:"tags,omitempty"`
	BusinessUnit   *ApplicationBusinessUnit `json:"business_unit,omitempty"`
	BusinessOwners []struct {
		Email string `json:"email,omitempty"`
		Name  string `json:"name,omitempty"`
	} `json:"business_owners,omitempty"`
	ArcherAppName       string              `json:"archer_app_name,omitempty"`
	Policies            []ApplicationPolicy `json:"policies,omitempty"`
	Teams               []ApplicationTeam   `json:"teams,omitempty"`
	CustomFields        []CustomField       `json:"custom_fields,omitempty"`
	Description         string              `json:"description,omitempty"`
	BusinessCriticality BusinessCriticality `json:"business_criticality,omitempty"`
}

type CustomField struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type ApplicationPolicy struct {
	Name      string `json:"name,omitempty"`
	Guid      string `json:"guid,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

type ApplicationTeam struct {
	Guid     string `json:"guid,omitempty"`
	TeamId   int    `json:"team_id,omitempty"`
	TeamName string `json:"team_name,omitempty"`
}

type ApplicationBusinessUnit struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Guid string `json:"guid,omitempty"`
}

type ApplicationBusinessOwner struct {
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

// ListApplicationOptions contains all of the fields that can be passed as query values when calling the ListApplications method.
// NOTE: the policy field is not currently included.
type ListApplicationOptions struct {
	Page                  int              `url:"page,omitempty"`
	Size                  int              `url:"size,omitempty"`
	Name                  string           `url:"name,omitempty"`                      // Filter Applications by Name (Not an exact match). Documentation Reference: https://docs.veracode.com/r/List_Applications_By_Name
	Tag                   string           `url:"tag,omitempty"`                       // Documentation Reference: https://docs.veracode.com/r/r_applications_any_tag and https://docs.veracode.com/r/r_applications_tag
	Team                  string           `url:"team,omitempty"`                      // Filter the Applications by team name.
	LegacyId              int              `url:"legacy_id,omitempty"`                 // Documentation Reference: https://docs.veracode.com/r/r_applications_info
	ScanType              ScanType         `url:"scan_type,omitempty"`                 // The valid scan_type values are STATIC, DYNAMIC and, for Manual Penetration Testing (MPT), MANUAL. Documentation Reference: https://docs.veracode.com/r/r_applications_scan_type
	ScanStatus            []ScanStatus     `url:"scan_status,omitempty"`               // Filter Applications by a list of scan statuses.
	BusinessUnit          string           `url:"business_unit,omitempty"`             // Return a list of Application Profiles that belong to the BU with this name. Documentation Reference: https://docs.veracode.com/r/r_applications_bu
	PolicyGuid            string           `url:"policy_guid,omitempty"`               // Filter Applications by the Policy that is assigned to them.
	PolicyCompliance      PolicyCompliance `url:"policy_compliance,omitempty"`         // Documentation Reference: https://docs.veracode.com/r/r_applications_compliance
	SortByCustomFieldName string           `url:"sort_by_custom_field_name,omitempty"` // Custom field name on which to sort.

	// You can use the Applications REST API to list the application profiles that have had an event that triggered a policy evaluation after a specific date.
	// The events that trigger policy evaluations are scans, approved mitigations, new component vulnerability releases, and policy changes.
	//
	// The value needs to be in format: 2006-01-02.
	//
	// Documentation Reference: https://docs.veracode.com/r/Listing_Applications_by_Last_Policy_Evaluation_Date_with_the_Applications_API
	PolicyComplianceCheckedAfter string `url:"policy_compliance_checked_after,omitempty"`

	// Send the following request to return the list of application profiles modified after a specific date.
	//
	// The value needs to be in format: 2006-01-02.
	//
	// Documentation Reference: https://docs.veracode.com/r/r_applications_modified_date
	ModifiedAfter string `url:"modified_after,omitempty"`

	// CustomFieldNames and CustomFieldValues need to both be set together.
	// You can use the AddCustomFieldOption method to set/update these fields.
	CustomFieldNames  []string `url:"custom_field_names,omitempty"`
	CustomFieldValues []string `url:"custom_field_values,omitempty"`
}

// AddCustomFieldOption sets the customFieldName and customFieldValue attributes on the ListApplicationOptions.
// To identify application profiles with any value for a specific custom field, enter the URL-encoded wildcard value %25 for customFieldValue.
//
// Documentation Reference: https://docs.veracode.com/r/r_applications_custom_field
func (l *ListApplicationOptions) AddCustomFieldOption(customFieldName, customFieldValue string) {
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

type applicationSearchResult struct {
	Embedded struct {
		Applications []Application `json:"applications"`
	} `json:"_embedded"`
	Links navLinks `json:"_links"`
	Page  pageMeta `json:"page"`
}

func (r *applicationSearchResult) GetLinks() navLinks {
	return r.Links
}

func (r *applicationSearchResult) GetPageMeta() pageMeta {
	return r.Page
}

// NewApplication creates an Application with all of the required fields.
func NewApplication(name, policyGuid string, businessCriticality BusinessCriticality) Application {
	return Application{
		Profile: ApplicationProfile{
			Name: name,
			Policies: []ApplicationPolicy{
				{Guid: policyGuid},
			},
			BusinessCriticality: businessCriticality,
		},
	}
}

// CreateApplication creates a new application using the provided Application.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/r_applications_create
//   - https://docs.veracode.com/r/r_applications_create_assign_team
func (a *ApplicationService) CreateApplication(ctx context.Context, application Application) (*Application, *Response, error) {
	byt, err := json.Marshal(application)
	if err != nil {
		return nil, nil, err
	}

	req, err := a.Client.NewRequest(ctx, "/appsec/v1/applications", http.MethodPost, bytes.NewBuffer(byt))
	if err != nil {
		return nil, nil, err
	}

	resp, err := a.Client.Do(req, &application)
	if err != nil {
		return nil, resp, err
	}
	return &application, resp, nil
}

// UpdateApplication updates the Application Profile provided.
// NOTE: When you update an application profile with this API, all properties are required.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/r_applications_update
//   - https://app.swaggerhub.com/apis/Veracode/veracode-applications_api_specification/1.0#/Application%20information%20API/updateApplicationUsingPUT
func (a *ApplicationService) UpdateApplication(ctx context.Context, application Application) (*Application, *Response, error) {
	buf, err := json.Marshal(application)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/applications/"+application.Guid, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var updatedApp Application
	resp, err := a.Client.Do(req, &updatedApp)
	if err != nil {
		return nil, resp, err
	}

	return &updatedApp, resp, err
}

// GetApplication retrieves an Application Profile with the provided appId.
//
// Veracode API documentation: https://app.swaggerhub.com/apis/Veracode/veracode-applications_api_specification/1.0#/Application%20information%20API/getApplicationUsingGET
func (a *ApplicationService) GetApplication(ctx context.Context, appId string) (*Application, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/applications/"+appId, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var getApp Application

	resp, err := a.Client.Do(req, &getApp)
	if err != nil {
		return nil, resp, err
	}

	return &getApp, resp, nil
}

// ListApplications takes a ListApplicationOptions and returns a list of Applications.
//
// Veracode API documentation: https://docs.veracode.com/r/r_applications_list
func (a *ApplicationService) ListApplications(ctx context.Context, options ListApplicationOptions) ([]Application, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/applications", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var result applicationSearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.Applications, resp, nil
}

// DeleteApplication deletes an application from the Veracode API using the provided appId.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-applications_api_specification/1.0#/Application%20information%20API/deleteApplicationUsingDELETE
//   - https://docs.veracode.com/r/r_applications_delete
func (a *ApplicationService) DeleteApplication(ctx context.Context, appId string) (*Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/applications/"+appId, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
