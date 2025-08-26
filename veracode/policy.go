package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PolicyCategory string
type PolicyScanType string
type PolicyScanFrequency string
type PolicyType string
type FindingRuleScanType string
type FindingRuleType string

const (
	CategoryApplication PolicyCategory = "APPLICATION"
	CategoryComponent   PolicyCategory = "COMPONENT"

	PolicyScanTypeStatic  PolicyScanType = "STATIC"
	PolicyScanTypeDynamic PolicyScanType = "DYNAMIC"
	PolicyScanTypeManual  PolicyScanType = "MANUAL"
	PolicyScanTypeSca     PolicyScanType = "SCA"
	PolicyScanTypeAny     PolicyScanType = "ANY"

	FindingRuleScanTypeStatic    FindingRuleScanType = "STATIC"
	FindingRuleScanTypeDynamic   FindingRuleScanType = "DYNAMIC"
	FindingRuleScanTypeManual    FindingRuleScanType = "MANUAL"
	FindingRuleScanTypesCA       FindingRuleScanType = "SCA"
	FindingRuleScanTypeMobile    FindingRuleScanType = "MOBILE"
	FindingRuleScanTypeAll       FindingRuleScanType = "ALL"
	FindingRuleScanTypeDynamicMP FindingRuleScanType = "DYNAMICMP"

	PolicyScanFrequencyNotRequired     PolicyScanFrequency = "NOT_REQUIRED"
	PolicyScanFrequencyOnce            PolicyScanFrequency = "ONCE"
	PolicyScanFrequencyWeekly          PolicyScanFrequency = "WEEKLY"
	PolicyScanFrequencyMonthly         PolicyScanFrequency = "MONTHLY"
	PolicyScanFrequencyQuarterly       PolicyScanFrequency = "QUARTERLY"
	PolicyScanFrequencySemiAnnually    PolicyScanFrequency = "SEMI_ANNUALLY"
	PolicyScanFrequencyAnnually        PolicyScanFrequency = "ANNUALLY"
	PolicyScanFrequencyEvery18Months   PolicyScanFrequency = "EVERY_18_MONTHS"
	PolicyScanFrequencyEvery2Years     PolicyScanFrequency = "EVERY_2_YEARS"
	PolicyScanFrequencyEvery3Years     PolicyScanFrequency = "EVERY_3_YEARS"
	PolicyScanFrequencySetByVLPolicy   PolicyScanFrequency = "SET_BY_VL_POLICY"
	PolicyScanFrequencySetByPolicyRule PolicyScanFrequency = "SET_BY_POLICY_RULE"

	PolicyTypeBuiltIn       PolicyType = "BUILTIN"
	PolicyTypeVeracodeLevel PolicyType = "VERACODELEVEL"
	PolicyTypeCustomer      PolicyType = "CUSTOMER"
	PolicyTypeStandard      PolicyType = "STANDARD"

	FindingRuleTypeFailAll          FindingRuleType = "FAIL_ALL"
	FindingRuleTypeCWE              FindingRuleType = "CWE"
	FindingRuleTypeCategory         FindingRuleType = "CATEGORY"
	FindingRuleTypeMaxSeverity      FindingRuleType = "MAX_SEVERITY"
	FindingRuleTypeCvss             FindingRuleType = "CVSS"
	FindingRuleTypeCve              FindingRuleType = "CVE"
	FindingRuleTypeBlacklist        FindingRuleType = "BLACKLIST"
	FindingRuleTypeMinScore         FindingRuleType = "MIN_SCORE"
	FindingRuleTypeSecurityStandard FindingRuleType = "SECURITY_STANDARD"
	FindingRuleTypeLicenseRisk      FindingRuleType = "LICENSE_RISK"
	FindingRuleTypeAllowlist        FindingRuleType = "ALLOWLIST"
)

// When creating a new Policy or creating a new version of an existing Policy, below fields are required:
//   - Name
//   - Description
//   - VendorPolicy
type PolicyVersion struct {
	Category                PolicyCategory   `json:"category,omitempty"`             // The category of the policy.
	Created                 *ctime           `json:"created"`                        // The date and time the application was created. The date/time is formatted as per RFC3339 and ISO-8601. The timezone is UTC.
	CustomSeverities        []CustomSeverity `json:"custom_severities,omitempty"`    // A set of severity overrides for use with this policy.
	Description             string           `json:"description,omitempty"`          // A description of the policy.
	EvaluationDate          *ctime           `json:"evaluation_date"`                // The date from which to define the evaluation timeframe, which determines when ﬁndings that violate rules should cause an application to not pass policy.
	EvaluationDateType      string           `json:"evaluation_date_type,omitempty"` // Specify the supported evaluation date type as before or after the specified date. [BEFORE, AFTER]
	FindingRules            []FindingRule    `json:"finding_rules,omitempty"`        // A set of rules to be evaluated against the scan findings.
	ScaGracePeriods         ScaGracePeriods  `json:"sca_grace_periods,omitempty"`
	Guid                    string           `json:"guid,omitempty"`                       // Unique identifier for the policy.
	ModifiedBy              string           `json:"modified_by,omitempty"`                // The username of the user who most recently modified the policy.
	Name                    string           `json:"name,omitempty"`                       // Policy name.
	OrganizationId          int              `json:"organization_id,omitempty"`            // The organization with which the policy is associated. If no value is provided, the results are publicly visible.
	ScaBlacklistGracePeriod int              `json:"sca_blacklist_grace_period,omitempty"` // (DEPRECATED) The grace period in number of days permitted for the component blocklist enforcement rule type.
	ScanFrequencyRules      []ScanFrequency  `json:"scan_frequency_rules,omitempty"`       // The set of scan frequencies to be evaluated.
	ScoreGracePeriod        int              `json:"score_grace_period,omitempty"`         // The number of days grace period allowed for the policy score.
	Sev0_grace_period       int              `json:"sev0_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 0.
	Sev1_grace_period       int              `json:"sev1_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 1.
	Sev2_grace_period       int              `json:"sev2_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 2.
	Sev3_grace_period       int              `json:"sev3_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 3.
	Sev4_grace_period       int              `json:"sev4_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 4.
	Sev5_grace_period       int              `json:"sev5_grace_period,omitempty"`          // The number of days grace period allowed for findings of severity 5.
	Type                    PolicyType       `json:"type,omitempty"`                       // The evaluation policy type.
	VendorPolicy            *bool            `json:"vendor_policy,omitempty"`              // Use this flag to indicate if this policy is to be visible and available for policy evaluation by a vendor organization.
	Version                 int              `json:"version,omitempty"`                    // The version of this policy.
}

type CustomSeverity struct {
	Cwe      int `json:"cwe,omitempty"`      // The CWE associated with the custom severity.
	Severity int `json:"severity,omitempty"` // The severity to be applied to findings of the specified CWE.
}

type ScanFrequency struct {
	Frequency     PolicyScanFrequency `json:"frequency,omitempty"`
	PolicyVersion PolicyVersion       `json:"policy_version,omitempty"`
	ScanType      PolicyScanType      `json:"scan_type,omitempty"`
}

type ScaGracePeriods struct {
	ScaBlacklistGracePeriod int                    `json:"sca_blacklist_grace_period,omitempty"` // The grace period in number of days permitted for the component blocklist enforcement rule type.
	LicenseRiskGracePeriod  int                    `json:"license_risk_grace_period,omitempty"`  // The grace period in number of days permitted for the component license risk rule type.
	SeverityGracePeriod     SeverityGracePeriod    `json:"severity_grace_period,omitempty"`
	CvssScoreGracePeriod    []CvssScoreGracePeriod `json:"cvss_score_grace_period,omitempty"` // The grace period in number of days permitted for the vulnerability CVSS score rule type.
}

type SeverityGracePeriod struct {
	Sev0GracePeriod int `json:"sev_0_grace_period,omitempty"`
	Sev1GracePeriod int `json:"sev_1_grace_period,omitempty"`
	Sev2GracePeriod int `json:"sev_2_grace_period,omitempty"`
	Sev3GracePeriod int `json:"sev_3_grace_period,omitempty"`
	Sev4GracePeriod int `json:"sev_4_grace_period,omitempty"`
}

type CvssScoreGracePeriod struct {
	Upper float64 `json:"upper,omitempty"` // The upper CVSS score limit for this grace period. Value must be between 0.0 and 10.0.
	Lower float64 `json:"lower,omitempty"` // The lower CVSS score limit for this grace period. Value must be between 0.0 and 10.0.
	Days  int     `json:"days,omitempty"`  // The grace period in number of days permitted for findings with a CVSS score within the range between the upper and lower CVSS score values.
}

type FindingRule struct {
	PolicyVersion   PolicyVersion              `json:"policy_version,omitempty"`
	Coordinate      Coordinate                 `json:"coordinate,omitempty"`
	ScanType        []FindingRuleScanType      `json:"scan_type,omitempty"` // The type of scan on which to enforce the rule.
	Type            FindingRuleType            `json:"type,omitempty"`      // Specify the supported rule types.
	Value           string                     `json:"value,omitempty"`     // The value of this specific rule, such as the minimal score value. This value does not apply to the FAIL_ALL rule type.
	AdvancedOptions FindingRuleAdvancedOptions `json:"advanced_options,omitempty"`
}

type Coordinate struct {
	Coordinate1 string `json:"coordinate_1,omitempty"` // The name of the first coordinate.
	Coordinate2 string `json:"coordinate_2,omitempty"` // The name of the second coordinate.
	CreatedBy   string `json:"created_by,omitempty"`   // The name of the user who created this coordinate.
	CreatedDate ctime  `json:"created_date,omitempty"` // The date when the user created the coordinate.
	FindingRule string `json:"finding_rule,omitempty"` //
	RepoType    string `json:"repo_type,omitempty"`    // The repository type of the coordinate; for example, nexus, or maven.
	Version     string `json:"version,omitempty"`      // The version of the coordinate.
}

type FindingRuleAdvancedOptions struct {
	AllLicensesMustMeetRequirement bool                `json:"all_licenses_must_meet_requirement,omitempty"` // Set to true to specify that component licenses must meet all policy rule requirements.
	AllowedNonossLicenses          bool                `json:"allowed_nonoss_licenses,omitempty"`            // Set to true to allow licenses that are not open-source (OSS).
	IsBlocklist                    bool                `json:"is_blocklist,omitempty"`                       // Set to true to add the selected list of licenses to the blocklist.
	SelectedLicenses               []ScaLicenseSummary `json:"selected_licenses,omitempty"`                  // List all selected licenses. (Must be unique)
	FindingRule                    string              `json:"finding_rule,omitempty"`
}

type ScaLicenseSummary struct {
	FullName string `json:"full_name,omitempty"` // Full name of the license.
	Name     string `json:"name,omitempty"`      // Short name of the license.
	Risk     string `json:"risk,omitempty"`      // Risk rating of the license. Values are Low, Medium, High, or Unknown.
	SpdxId   string `json:"spdx_id,omitempty"`   // SPDX identifier for the license.
	Url      string `json:"url,omitempty"`       // URL to the license on the spdx.org website.
}

type ListPolicyOptions struct {
	Category       PolicyCategory `url:"category,omitempty"`         // The category of the policy. [APPLICATION, COMPONENT]
	LegacyPolicyId int            `url:"legacy_policy_id,omitempty"` // Filters results based on the ID of the custom policy created in the Veracode Platform.
	Name           string         `url:"name,omitempty"`             // Filter on the policy name.
	NameExact      string         `url:"name_exact,omitempty"`       // Use this flag to enforce exact name-matching when filtering on the policy name.
	Page           int            `url:"page,omitempty"`             // Page number. Defaults to 0.
	PublicPolicy   *bool          `url:"public_policy,omitempty"`    // Filters results to include or exclude a public Veracode policy. Default: true
	Size           int            `url:"size,omitempty"`             // Page size (1-500, defaults to 50).
	VendorPolicy   bool           `url:"vendor_policy,omitempty"`    // Filters results to those with or without a vendor policy flag.
}

type ListPolicyVersionsOptions struct {
	Page int `url:"page,omitempty"` // Page number. Defaults to 0.
	Size int `url:"size,omitempty"` // Page size (1-500, defaults to 50).
}

type policySearchResult struct {
	Embedded struct {
		PolicyVersions []PolicyVersion `json:"policy_versions"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *policySearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *policySearchResult) GetPageMeta() PageMeta {
	return r.Page
}

type scaLicenseSummaryListResult struct {
	Embedded struct {
		ScaLicenseSummary []ScaLicenseSummary `json:"sca_license_summaries"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *scaLicenseSummaryListResult) GetLinks() NavLinks {
	return r.Links
}

func (r *scaLicenseSummaryListResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListPolicies takes a [ListPolicyOptions] and returns a list of [PolicyVersion].
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/PolicyVersion
func (a *PolicyService) ListPolicies(ctx context.Context, options ListPolicyOptions) ([]PolicyVersion, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policies", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var result policySearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.PolicyVersions, resp, nil
}

// GetPolicy retrieves a [PolicyVersion] with the provided policyGuid.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) GetPolicy(ctx context.Context, policyGuid string) (*PolicyVersion, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policies/"+policyGuid, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var getPolicy PolicyVersion

	resp, err := a.Client.Do(req, &getPolicy)
	if err != nil {
		return nil, resp, err
	}

	return &getPolicy, resp, nil
}

// ListPolicies takes a [ListPolicyVersionsOptions] and a policyGuid, and returns a list of [PolicyVersion] for the give Policy.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) ListPolicyVersions(ctx context.Context, policyGuid string, options ListPolicyVersionsOptions) ([]PolicyVersion, *Response, error) {
	req, err := a.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/policies/%s/versions", policyGuid), http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var result policySearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.PolicyVersions, resp, nil
}

// GetPolicyVersion retrieves a specific [PolicyVersion] by version.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) GetPolicyVersion(ctx context.Context, policyGuid string, version int) (*PolicyVersion, *Response, error) {
	req, err := a.Client.NewRequest(ctx, fmt.Sprintf("/appsec/v1/policies/%s/versions/%d", policyGuid, version), http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var getPolicy PolicyVersion

	resp, err := a.Client.Do(req, &getPolicy)
	if err != nil {
		return nil, resp, err
	}

	return &getPolicy, resp, nil
}

// UpdatePolicy creates a new Policy Version that will become the default when making a get request.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) UpdatePolicy(ctx context.Context, policyGuid string, policyVersion PolicyVersion) (*PolicyVersion, *Response, error) {
	buf, err := json.Marshal(policyVersion)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policies/"+policyGuid, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var updatedPolicy PolicyVersion
	resp, err := a.Client.Do(req, &updatedPolicy)
	if err != nil {
		return nil, resp, err
	}

	return &updatedPolicy, resp, err
}

// CreatePolicy creates a new Policy with version 1.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) CreatePolicy(ctx context.Context, policyVersion PolicyVersion) (*PolicyVersion, *Response, error) {
	buf, err := json.Marshal(policyVersion)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policies", http.MethodPost, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var createdPolicy PolicyVersion
	resp, err := a.Client.Do(req, &createdPolicy)
	if err != nil {
		return nil, resp, err
	}

	return &createdPolicy, resp, err
}

// DeletePolicy deletes a Policy with the provided policyGuid.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) DeletePolicy(ctx context.Context, policyGuid string) (*Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policies/"+policyGuid, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// ListSCAComponentLicenseInformation returns a list component licenses associated with SCA security policies.
//
// API Documentation: https://app.swaggerhub.com/apis/Veracode/veracode-policy_api_specification/1.0#/
func (a *PolicyService) ListSCAComponentLicenseInformation(ctx context.Context, options PageOptions) ([]ScaLicenseSummary, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_licenselist", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var result scaLicenseSummaryListResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.ScaLicenseSummary, resp, nil
}
