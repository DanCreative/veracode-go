package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type ComponentPolicySetting struct {
	Modified   *ctime `json:"modified,omitempty"`    // The date and time when the pre-build component default policy setting was modified. The date and time format is per RFC3339 and ISO-8601. Timezone is UTC.
	ModifiedBy string `json:"modified_by,omitempty"` // Name of the user who most recently modified the pre-build component default policy setting.
	Plugin     string `json:"plugin,omitempty"`
	PolicyGuid string `json:"policy_guid,omitempty"` // Unique identifier for the pre-build component policy.
}

type PolicySetting struct {
	BusinessCriticality BusinessCriticality `json:"business_criticality,omitempty"` // The business criticality for which this policy is the default policy.
	Modified            *ctime              `json:"modified,omitempty"`             // The date and time when the application default policy setting was modified. The date and time format is per RFC3339 and ISO-8601. Timezone is UTC.
	PolicyGuid          string              `json:"policy_guid,omitempty"`          // Unique identifier for the application policy.
	AgentSetting        string              `json:"agent_setting,omitempty"`
}

type policySettingSearchResult struct {
	Embedded struct {
		PolicySettings []PolicySetting `json:"policy_settings"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *policySettingSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *policySettingSearchResult) GetPageMeta() PageMeta {
	return r.Page
}

type componentPolicySearchResult struct {
	Embedded struct {
		ComponentPolicySettings []ComponentPolicySetting `json:"component_policy_settings"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *componentPolicySearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *componentPolicySearchResult) GetPageMeta() PageMeta {
	return r.Page
}

// ListPolicySettings returns a slice of [PolicySetting] for the application policies.
func (a *PolicyService) ListPolicySettings(ctx context.Context) ([]PolicySetting, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var result policySettingSearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.PolicySettings, resp, nil
}

// UpdatePolicySettings takes a slice of [PolicySetting] and updates all of the settings present in the slice.
func (a *PolicyService) UpdatePolicySettings(ctx context.Context, policySettings []PolicySetting) ([]PolicySetting, *Response, error) {
	buf, err := json.Marshal(policySettings)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings", http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var updatedPolicySettings []PolicySetting
	resp, err := a.Client.Do(req, updatedPolicySettings)
	if err != nil {
		return nil, resp, err
	}

	return updatedPolicySettings, resp, err
}

// ListComponentPolicySettings returns the default pre-build component policies.
func (a *PolicyService) ListComponentPolicySettings(ctx context.Context) ([]ComponentPolicySetting, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings/components", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var result componentPolicySearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.ComponentPolicySettings, resp, nil
}

// UpdateComponentPolicySetting updates the pre-build component default policies of an organization.
func (a *PolicyService) UpdateComponentPolicySetting(ctx context.Context, componentPolicySettings []ComponentPolicySetting) ([]ComponentPolicySetting, *Response, error) {
	buf, err := json.Marshal(componentPolicySettings)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings/components", http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var updatedComponentPolicySettings []ComponentPolicySetting
	resp, err := a.Client.Do(req, &updatedComponentPolicySettings)
	if err != nil {
		return nil, resp, err
	}

	return updatedComponentPolicySettings, resp, err
}

// ListThirdPartyComponentPolicySettings returns the default pre-build component policies.
func (a *PolicyService) ListThirdPartyComponentPolicySettings(ctx context.Context) ([]ComponentPolicySetting, *Response, error) {
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings/thirdparty_libraries", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var result componentPolicySearchResult

	resp, err := a.Client.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.ComponentPolicySettings, resp, nil
}

// UpdateThirdPartyComponentPolicySettings updates the pre-build component default policies of an organization.
func (a *PolicyService) UpdateThirdPartyComponentPolicySettings(ctx context.Context, componentPolicySettings []ComponentPolicySetting) ([]ComponentPolicySetting, *Response, error) {
	buf, err := json.Marshal(componentPolicySettings)
	if err != nil {
		return nil, nil, err
	}
	req, err := a.Client.NewRequest(ctx, "/appsec/v1/policy_settings/components", http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var updatedComponentPolicySettings []ComponentPolicySetting
	resp, err := a.Client.Do(req, &updatedComponentPolicySettings)
	if err != nil {
		return nil, resp, err
	}

	return updatedComponentPolicySettings, resp, err
}
