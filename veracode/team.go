package veracode

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Team struct {
	TeamId       string           `json:"team_id,omitempty"`
	TeamLegacyId int              `json:"team_legacy_id,omitempty"`
	TeamName     string           `json:"team_name,omitempty"`
	Relationship TeamRelationship `json:"relationship,omitempty"`
}

type TeamRelationship struct {
	Name string `json:"name,omitempty"`
}

// teamSearchResult is required to decode the list of teams and search user response bodies.
type teamSearchResult struct {
	Embedded struct {
		Teams []Team `json:"teams"`
	} `json:"_embedded"`
	Links navLinks `json:"_links"`
	Page  pageMeta `json:"page"`
}

// ListTeamOptions contains all of the fields that can be passed as query values.
type ListTeamOptions struct {
	Size      int   `url:"size,omitempty"`
	Page      int   `url:"page"`
	AllForOrg *bool `url:"all_for_org"`
}

// If Relationship.Name is "", create custom struct where TeamRelationship is a pointer and set it to nil.
// This will omit relationship from the marshalled json.
//
// If Relationship.Name is not "", flatten TeamRelationship to Relationship in Team model.
func (t *Team) MarshalJSON() ([]byte, error) {
	type Alias Team
	if t.Relationship.Name == "" {
		return json.Marshal(&struct {
			*Alias
			Relationship *TeamRelationship `json:"relationship,omitempty"`
		}{
			Alias:        (*Alias)(t),
			Relationship: nil,
		})
	}
	return json.Marshal(&struct {
		*Alias
		Relationship string `json:"relationship,omitempty"`
	}{
		Alias:        (*Alias)(t),
		Relationship: t.Relationship.Name,
	})
}

// ListTeams takes a ListTeamsOptions and returns a list of teams.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_list_teams
func (i *IdentityService) ListTeams(ctx context.Context, options ListTeamOptions) ([]Team, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/teams", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = values.Encode()

	var teamsResult teamSearchResult

	resp, err := i.Client.Do(req, &teamsResult)
	if err != nil {
		return nil, resp, err
	}
	return teamsResult.Embedded.Teams, resp, err
}
