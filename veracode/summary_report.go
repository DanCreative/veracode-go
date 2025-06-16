package veracode

type SummaryReportOptions struct {
	BuildId int    `url:"build_id,omitempty"` // ID of the build in which the scan ran. Default is the latest build_id.
	Context string `url:"context,omitempty"`  // GUID of the associated development sandbox, if specified. The Summary Report is relative to this context parameter.
}

type SummaryReport struct {
	StaticAnalysis              AnalysisType                `json:"static-analysis"`
	DynamicAnalysis             AnalysisType                `json:"dynamic-analysis"`
	ManualAnalysis              ManualAnalysisType          `json:"manual-analysis"`
	Severity                    []SeverityType              `json:"severity"`
	FlawStatus                  FlawStatusType              `json:"flaw_status"`
	CustomFields                CustomFields                `json:"custom_fields"`
	SoftwareCompositionAnalysis SoftwareCompositionAnalysis `json:"software_composition_analysis"`
	ReportFormatVersion         string                      `json:"report_format_version,omitempty"`    // Version of the format of this report.
	AccountId                   int                         `json:"account_id,omitempty"`               // ID of the Veracode account.
	AppName                     string                      `json:"app_name,omitempty"`                 // Name of the scanned application.
	AppId                       int                         `json:"app_id,omitempty"`                   // ID of the scanned application.
	AnalysisId                  int                         `json:"analysis_id,omitempty"`              // ID for the scan.
	StaticAnalysisUnitId        int                         `json:"static_analysis_unit_id,omitempty"`  // Unit ID for a static analysis.
	SandboxName                 string                      `json:"sandbox_name,omitempty"`             // Name of the development sandbox. Not applicable for a policy scan.
	SandboxId                   int                         `json:"sandbox_id,omitempty"`               // ID of the development sandbox. Not applicable for a policy scan.
	FirstBuildSubmittedDate     ctime                       `json:"first_build_submitted_date"`         // Timestamp of the first time you submitted a build of this application to Veracode for scanning.
	Version                     string                      `json:"version,omitempty"`                  // Version label for the application.
	BuildId                     int                         `json:"build_id,omitempty"`                 // ID of the build for the application.
	Vendor                      string                      `json:"vendor,omitempty"`                   // Name of the vendor that provided the application, if applicable.
	Submitter                   string                      `json:"submitter,omitempty"`                // Name of the account or user that created the build.
	Platform                    string                      `json:"platform,omitempty"`                 // Platform of the build for the application.
	BusinessCriticality         int                         `json:"business_criticality,omitempty"`     // Business criticality for the application.
	GenerationDate              ctime                       `json:"generation_date"`                    // Timestamp when Veracode generated the report.
	VeracodeLevel               string                      `json:"veracode_level,omitempty"`           // Security score for the application based on Veracode Levels. Values are VL1, VL2, VL3, VL4, or VL5
	TotalFlaws                  int                         `json:"total_flaws,omitempty"`              // Total number of discovered findings for the application.
	FlawsNotMitigated           int                         `json:"flaws_not_mitigated,omitempty"`      // Total number of discovered findings not marked as mitigated.
	Teams                       string                      `json:"teams,omitempty"`                    // Teams assigned to this application.
	LifeCycleStage              string                      `json:"life_cycle_stage,omitempty"`         // Current life cycle stage for this application. For example, deployed or in development.
	PlannedDeploymentDate       ctime                       `json:"planned_deployment_date"`            // Deployment date for the application, if specified.
	LastUpdateTime              ctime                       `json:"last_update_time"`                   // Last time this application was modified.
	IsLatestBuild               bool                        `json:"is_latest_build,omitempty"`          // True if this report is for the most recent build of this application.
	PolicyName                  string                      `json:"policy_name,omitempty"`              // Name of the security policy assigned to this application.
	PolicyVersion               int                         `json:"policy_version,omitempty"`           // Version number of the security policy assigned to the version of this application.
	PolicyComplianceStatus      string                      `json:"policy_compliance_status,omitempty"` // Current policy compliance status for this application. Values are Calculating, Did Not Pass, Conditional Pass, or Pass.
	PolicyRulesStatus           string                      `json:"policy_rules_status,omitempty"`      // Current policy rules compliance status for this application. Does not include scan frequency requirements and grace period time allowed to address rule violations. Values are Calculating, Did Not Pass, or Pass.
	GracePeriodExpired          bool                        `json:"grace_period_expired,omitempty"`     // True if findings in the latest analyzed build of this application have existed for longer than the allowed grace period.
	ScanOverdue                 string                      `json:"scan_overdue,omitempty"`             // True if the amount of time between the last analysis and the current time is greater than the scan frequency that your security policy requires.
	AnyTypeScanDue              ctime                       `json:"any_type_scan_due"`                  // Date to analyze a new build of this application for it to remain in compliance with the required scan frequency of the security policy.
	BusinessOwner               string                      `json:"business_owner,omitempty"`           // First and last name of the party responsible for this application.
	BusinessUnit                string                      `json:"business_unit,omitempty"`            // Department or group associated with this application.
	Tags                        string                      `json:"tags,omitempty"`                     // Comma-delimited list of tags associated with this application.
	LegacyScanEngine            bool                        `json:"legacy_scan_engine,omitempty"`       // For a static analysis, indicates whether the scan ran with a legacy engine or the same engine version as the previous scan of its type.
}

// For a static analysis, a list of modules with one module node per module analyzed. For a dynamic analysis, a single module node.
type AnalysisType struct {
	Modules            Module `json:"modules"`
	Rating             string `json:"rating,omitempty"`                // Letter grade for the security of this application.
	Score              int    `json:"score,omitempty"`                 // Numeric security score for this application.
	MitigatedRating    string `json:"mitigated_rating,omitempty"`      // Letter grade for the security of this application, based on mitigated findings.
	MitigatedScore     int    `json:"mitigated_score,omitempty"`       // Numeric security score for this application, based on mitigated findings.
	SubmittedDate      ctime  `json:"submitted_date"`                  // Date when you submitted this application to Veracode for analysis.
	PublishedDate      ctime  `json:"published_date"`                  // Date when Veracode published the analysis for this application.
	NextScanDue        ctime  `json:"next_scan_due"`                   // Date when the active security policy for this application is scheduled to request the next scan.
	AnalysisSizeBytes  int    `json:"analysis_size_bytes,omitempty"`   // Optional. For a static analysis, the size, in bytes, of the scanned modules.
	EngineVersion      string `json:"engine_version,omitempty"`        // For a static analysis, the version of the engine that Veracode used for this scan.
	DynamicScanType    string `json:"dynamic_scan_type,omitempty"`     // Optional. For a dynamic analysis, indicates whether the scan is DA (Dynamic Analysis), MP (DynamicMP), or DS (DynamicDS).
	ScanExitStatusId   int    `json:"scan_exit_status_id,omitempty"`   // Optional. For a dynamic analysis, the numeric code for scan exit status.
	ScanExitStatusDesc string `json:"scan_exit_status_desc,omitempty"` // Optional. For a dynamic analysis, a description for scan_exit_status_id.
	Version            string `json:"version,omitempty"`               // Optional. Version of the scan.
}

type Module struct {
	Module []ModuleType `json:"module,omitempty"`
}

// Information about the type of module that Veracode scanned.
type ModuleType struct {
	Name         string `json:"name,omitempty"`            // Name of the scanned module. For a dynamic analysis, the name is blank.
	Compiler     string `json:"compiler,omitempty"`        // Compiler that compiled the scanned module. For a dynamic analysis, the value is blank.
	Os           string `json:"os,omitempty"`              // Operating system for which the scanned module is targetted. For a dynamic analysis, the value is blank.
	Architecture string `json:"architecture,omitempty"`    // Target architecture for which the scanned module is targeted. For a dynamic analysis, the value is blank.
	Loc          int    `json:"loc,omitempty"`             // Lines of codes. For a dynamic analysis or non-debug modules, the value is blank.
	Score        int    `json:"score,omitempty"`           // Module-specific security score, which contributes toward the analysis scores for the application.
	NumFlawsSev0 int    `json:"num_flaws_sev_0,omitempty"` // Number of severity-0 findings. These findings are the lowest severity and are usually informational only.
	NumFlawsSev1 int    `json:"num_flaws_sev_1,omitempty"` // Number of severity-1 findings.
	NumFlawsSev2 int    `json:"num_flaws_sev_2,omitempty"` // Number of severity-2 findings.
	NumFlawsSev3 int    `json:"num_flaws_sev_3,omitempty"` // Number of severity-3 findings.
	NumFlawsSev4 int    `json:"num_flaws_sev_4,omitempty"` // Number of severity-4 findings.
	NumFlawsSev5 int    `json:"num_flaws_sev_5,omitempty"` // Number of severity-5 findings. These findings are the highest severity and Veracode recommends that you fix them immediately.
	TargetUrl    string `json:"target_url,omitempty"`      // For a dynamic analysis, the URL for the application you scanned.
	Domain       string `json:"domain,omitempty"`          // For a dynamic analysis, the domain for the application you scanned.
}

// For Manual Penetration Testing, Veracode applies the confidentiality, integrity, and availability (CIA) triad to generate the final numeric score for the application.
// The report lists the delivery consultants, if any, followed by scan results.
type ManualAnalysisType struct {
	CiaAdjustment      int      `json:"cia_adjustment,omitempty"`      // For Manual Penetration Testing, the CIA triad that Veracode applied to the security score.
	Rating             string   `json:"rating,omitempty"`              // Letter grade for the security of this application.
	Score              int      `json:"score,omitempty"`               // Numeric score for the security of this application.
	NextScanDue        ctime    `json:"next_scan_due"`                 // Date when the active security policy for this application is scheduled to request the next scan.
	DeliveryConsultant []string `json:"delivery_consultant,omitempty"` // For Manual Penetration Testing, the names of the delivery consultants, if any.
	Modules            Module   `json:"modules"`
}

// Information about the Veracode Levels for the severity of a finding. The range is 0 through 5, where 0 is informational and 5 is the most severe.
type SeverityType struct {
	Level    int            `json:"level,omitempty"` // Veracode Level for the severity of the finding. The value range is 0 to 5, with 5 being the highest severity.
	Category []CategoryType `json:"category,omitempty"`
}

type CategoryType struct {
	CategoryName string `json:"category_name,omitempty"` // Name of the severity category.
	Severity     string `json:"severity,omitempty"`      // Enum: Informational, Very Low, Low, Medium, High, Very High
	Count        int    `json:"count,omitempty"`         // Number of findings in this category.
}

type FlawStatusType struct {
	New                      int `json:"new,omitempty"`                        // Number of findings discovered during the first build of this application.
	Reopen                   int `json:"reopen,omitempty"`                     // Number of findings discovered in a prior build of this application that were not new, but Veracode discovered them in the build immediately prior to this build.
	Open                     int `json:"open,omitempty"`                       // Number of findings discovered in this build that Veracode also discovered in the build immediately prior to this build.
	Fixed                    int `json:"fixed,omitempty"`                      // Number of findings discovered in the prior build that Veracode did not discover in the current build. For a dyanamic analysis, Veracode verifies the findings as fixed.
	Total                    int `json:"total,omitempty"`                      // Total number of findings discovered in this build.
	NotMitigated             int `json:"not_mitigated,omitempty"`              // Total number of findings discovered in this build that are not mitigated.
	Sev1Change               int `json:"sev-1-change,omitempty"`               // Number of severity-1 findings discovered in this build, minus the number of severity-1 findings discovered in the build immediately prior to this build.
	Sev2Change               int `json:"sev-2-change,omitempty"`               // Number of severity-2 findings discvoered in this build, minus the number of severity-2 findings discovered in the build immediately prior to this build.
	Sev3Change               int `json:"sev-3-change,omitempty"`               // Number of severity-3 findings discvoered in this build, minus the number of severity-3 findings discovered in the build immediately prior to this build.
	Sev4Change               int `json:"sev-4-change,omitempty"`               // Number of severity-4 findings discvoered in this build, minus the number of severity-4 findings discovered in the build immediately prior to this build.
	Sev5Change               int `json:"sev-5-change,omitempty"`               // Number of severity-5 findings discvoered in this build, minus the number of severity-5 findings discovered in the build immediately prior to this build.
	ConformsToGuidelines     int `json:"conforms-to-guidelines,omitempty"`     // Number of mitigations that adhere to your risk tolerance guidelines based on Veracode review.
	DeviatesFromGuidelines   int `json:"deviates-from-guidelines,omitempty"`   // Number of mitigations that either do not provide enough information or do not adhere to your the risk tolerance guidelines, based on Veracode review.
	TotalReviewedMitigations int `json:"total-reviewed-mitigations,omitempty"` // Total number of mitigations that Veracode reviewed. The value may not add up to the total number of all proposed or accepted mitigations.
}

type CustomFields struct {
	CustomField []CustomField `json:"custom_field,omitempty"`
}

// Information about findings discovered during Software Composition Analysis (SCA).
type SoftwareCompositionAnalysis struct {
	VulnerableComponents     VulnerableComponentList `json:"vulnerable_components,omitempty"`
	ThirdPartyComponents     int                     `json:"third_party_components,omitempty"`     // Number of vulnerable third party components.
	ViolatePolicy            bool                    `json:"violate_policy,omitempty"`             // Whether the component violates the security policy.
	ComponentsViolatedPolicy int                     `json:"components_violated_policy,omitempty"` // Number of components that violate the SCA policy.
	BlacklistedComponents    int                     `json:"blacklisted_components,omitempty"`     // Number of blacklisted components.
	ScaServiceAvailable      bool                    `json:"sca_service_available,omitempty"`      // True if the SCA service is available, else set to false.
}

// Details about the vulnerable components.
type VulnerableComponentList struct {
	Component []Component `json:"component_dto,omitempty"`
}

type Component struct {
	ComponentId                      string            `json:"component_id,omitempty"`                        // ID of the component.
	FileName                         string            `json:"file_name,omitempty"`                           // Filename of the component.
	Sha1                             string            `json:"sha1,omitempty"`                                // sha1
	Vulnerability                    int               `json:"vulnerability,omitempty"`                       // Number of vulnerabilities that Veracode discovered in the component.
	MaxCvssScore                     string            `json:"max_cvss_score,omitempty"`                      // Max Common Vulnerability Scoring System (CVSS) of the component. See cvss_score.
	Library                          string            `json:"library,omitempty"`                             // Library name of the component.
	Version                          string            `json:"version,omitempty"`                             // Version of the component.
	Vendor                           string            `json:"vendor,omitempty"`                              // Vendor name of the component.
	Description                      string            `json:"description,omitempty"`                         // Description of the component.
	Blacklisted                      string            `json:"blacklisted,omitempty"`                         // Blacklisted status for the component.
	New                              string            `json:"new,omitempty"`                                 // Whether this is a newly-added component.
	AddedDate                        ctime             `json:"added_date"`                                    // Date when you added the component.
	ComponentAffectsPolicyCompliance string            `json:"component_affects_policy_compliance,omitempty"` // Whether the component violates the SCA policy.
	FilePaths                        FilePathList      `json:"file_paths"`
	LicenseList                      LicenseList       `json:"licenses"`
	Vulnerabilities                  VulnerabilityList `json:"vulnerabilities"`
	ViolatedPolicyRules              ViolatedRuleList  `json:"violated_policy_rules"`
}

// Filepaths for the component.
type FilePathList struct {
	FilePath []FilePath `json:"file_path,omitempty"`
}

type FilePath struct {
	Value string `json:"value,omitempty"` // Filepath for the component.
}

// License details for the component.
type LicenseList struct {
	Licenses []License `json:"license_dto,omitempty"`
}

type License struct {
	Name       string `json:"name,omitempty"`        // Name of this license.
	SpdxId     string `json:"spdx_id,omitempty"`     // Classification for the license from the Software Package Data Exchange (SPDX) license list.
	LicenseUrl string `json:"license_url,omitempty"` // URL for this license.
	RiskRating string `json:"risk_rating,omitempty"` // Risk associated with the use of this license.
}

type VulnerabilityList struct {
	Vulnerability []Vulnerability `json:"vulnerability_dto"`
}

type Vulnerability struct {
	CveId                                string  `json:"cve_id,omitempty"`                                  // Common Vulnerabilities and Exposures (CVE) ID of the vulnerability.
	CvssScore                            float32 `json:"cvss_score,omitempty"`                              // Common Vulnerability Scoring System (CVSS) score. Measures the level of complexity for the vulnerability. The value is a range of 0 to 10 with 10 representing the highest complexity.
	Severity                             int     `json:"severity,omitempty"`                                // Veracode Level for the severity of the vulnerability. The value range is 0 to 5, with 5 being the highest severity.
	CweId                                string  `json:"cwe_id,omitempty"`                                  // Common Weakness Enumration (CWE) ID for the vulnerability.
	FirstFoundDate                       ctime   `json:"first_found_date"`                                  // Date when Veracode first discovered the vulnerability.
	CweSummary                           string  `json:"cwe_summary,omitempty"`                             // CVE summary for the vulnerability.
	SeverityDesc                         string  `json:"severity_desc,omitempty"`                           // Severity description for the vulnerbseverity.
	Mitigation                           string  `json:"mitigation,omitempty"`                              // Vulnerability mitigation status.
	MitigationType                       string  `json:"mitigation_type,omitempty"`                         // Type of mitigation applied to the vulnerability, if any.
	MitigatedDate                        ctime   `json:"mitigated_date"`                                    // Mitigation date for teh vulnerability.
	VulnerabilityAffectsPolicyCompliance string  `json:"vulnerability_affects_policy_compliance,omitempty"` // Whether the vulnerability affects SCA policy compliance.
}

type ViolatedRuleList struct {
	PolicyRule []PolicyRule `json:"policy_rule"`
}

type PolicyRule struct {
	Type  string `json:"type,omitempty"`  // Enum: DISALLOW_VULNERABILITIES_BY_SEVERITY, DISALLOW_CVSS_SCORE, DISALLOW_COMPONENT_BLACKLIST, DISALLOW_COMPONENT_BY_LICENSE_RISK
	Value string `json:"value,omitempty"` // SCA policy type.
	Desc  string `json:"desc,omitempty"`  // SCA policy description.
}
