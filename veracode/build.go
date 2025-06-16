package veracode

import (
	"context"
	"encoding/xml"
	"net/http"
	"time"
)

type BuildList struct {
	XMLName          xml.Name       `xml:"buildlist"`
	BuildListVersion string         `xml:"buildlist_version,attr"`
	AccountId        string         `xml:"account_id,attr"`
	AppId            string         `xml:"app_id,attr"`
	AppName          string         `xml:"app_name,attr"`
	Builds           []BuildSummary `xml:"build"`
}

type BuildInfo struct {
	XMLName          xml.Name      `xml:"buildinfo"`
	BuildInfoVersion string        `xml:"buildinfo_version,attr"`
	AccountId        string        `xml:"account_id,attr"`
	AppId            string        `xml:"app_id,attr"`
	BuildId          string        `xml:"build_id,attr"`
	Build            BuildDetailed `xml:"build"`
}

type BuildSummary struct {
	BuildId           string    `xml:"build_id,attr"`
	Version           string    `xml:"version,attr"`
	PolicyUpdatedDate time.Time `xml:"policy_updated_date,attr"`
	DynamicScanType   string    `xml:"dynamic_scan_type,attr"`
}

type BuildDetailed struct {
	XMLName                xml.Name     `xml:"build"`
	Version                string       `xml:"version,attr"`
	BuildId                string       `xml:"build_id,attr"`
	Submitter              string       `xml:"submitter,attr"`
	Platform               string       `xml:"platform,attr"`
	LifeCycleStage         string       `xml:"lifecycle_stage,attr"`
	SCAResultsReady        bool         `xml:"sca_results_ready,attr"`
	ResultsReady           bool         `xml:"results_ready,attr"`
	PolicyName             string       `xml:"policy_name,attr"`
	PolicyVersion          string       `xml:"policy_version,attr"`
	PolicyComplianceStatus string       `xml:"policy_compliance_status,attr"`
	PolicyUpdatedDate      time.Time    `xml:"policy_updated_date,attr"`
	RulesStatus            string       `xml:"rules_status,attr"`
	GracePeriodExpired     bool         `xml:"grace_period_expired,attr"`
	ScanOverdue            bool         `xml:"scan_overdue,attr"`
	LegacyScanEngine       bool         `xml:"legacy_scan_engine,attr"`
	AnalysisUnit           AnalysisUnit `xml:"analysis_unit"`
}

type AnalysisUnit struct {
	AnalysisType         string    `xml:"analysis_type,attr"`
	PublishedDate        time.Time `xml:"published_date,attr"`
	Status               string    `xml:"status,attr"`
	PublishedDateSeconds int       `xml:"published_date_sec,attr"`
	EngineVersion        int       `xml:"engine_version,attr"`
}

type BuildInfoOptions struct {
	AppId     int `url:"app_id,omitempty"`     // AppId is required
	BuildId   int `url:"build_id,omitempty"`   // Application or sandbox build ID. Default is the most recent static scan
	SandboxId int `url:"sandbox_id,omitempty"` // Target Sandbox Id
}

type BuildListOptions struct {
	AppId     int `url:"app_id,omitempty"`     // AppId is required
	SandboxId int `url:"sandbox_id,omitempty"` // Target Sandbox Id
}

// GetBuildInfo provides information about the most recent scan or a specific scan of the application.
//
// Documentation Reference: https://docs.veracode.com/r/r_getbuildinfo
func (u *UploadXMLService) GetBuildInfo(ctx context.Context, options BuildInfoOptions) (BuildInfo, *Response, error) {
	req, err := u.Client.NewRequest(ctx, "/api/5.0/getbuildinfo.do", http.MethodGet, nil, true)
	if err != nil {
		return BuildInfo{}, nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	req.URL.RawQuery = QueryEncode(options)

	var result BuildInfo

	resp, err := u.Client.Do(req, &result)
	if err != nil {
		return BuildInfo{}, resp, err
	}

	return result, resp, nil
}

func (u *UploadXMLService) GetBuildList(ctx context.Context, options BuildListOptions) (BuildList, *Response, error) {
	req, err := u.Client.NewRequest(ctx, "/api/5.0/getbuildlist.do", http.MethodGet, nil, true)
	if err != nil {
		return BuildList{}, nil, err
	}

	req.Header.Set("Content-Type", "application/xml")

	req.URL.RawQuery = QueryEncode(options)

	var result BuildList

	resp, err := u.Client.Do(req, &result)
	if err != nil {
		return BuildList{}, resp, err
	}

	return result, resp, nil
}
