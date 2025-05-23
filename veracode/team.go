package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Team struct {
	TeamId       string           `json:"team_id,omitempty"`
	TeamLegacyId int              `json:"team_legacy_id,omitempty"`
	TeamName     string           `json:"team_name,omitempty"`
	Relationship TeamRelationship `json:"relationship,omitempty"`
	Users        *[]User          `json:"users,omitempty"`
	BusinessUnit *BusinessUnit    `json:"business_unit,omitempty"`
}

type TeamRelationship struct {
	Name string `json:"name,omitempty"`
}

// ListTeamOptions contains all of the fields that can be passed as query values.
type ListTeamOptions struct {
	AllForOrg       *bool  `url:"all_for_org,omitempty"`
	TeamName        string `url:"team_name,omitempty"`
	IgnoreSelfTeams *bool  `url:"ignore_self_teams,omitempty"` // If true, return all teams in the organization. If false, return the teams the current user is a part of.
	OnlyManageable  bool   `url:"only_manageable,omitempty"`   // Only return teams manageable by the requesting user.
	Deleted         bool   `url:"deleted,omitempty"`           // Returns deleted teams.
	PageOptions            // can sort team_name field
}

// teamSearchResult is required to decode the list of teams and search user response bodies.
type teamSearchResult struct {
	Embedded struct {
		Teams []Team `json:"teams"`
	} `json:"_embedded"`
	Links NavLinks `json:"_links"`
	Page  PageMeta `json:"page"`
}

func (r *teamSearchResult) GetLinks() NavLinks {
	return r.Links
}

func (r *teamSearchResult) GetPageMeta() PageMeta {
	return r.Page
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
func (i *IdentityService) ListTeams(ctx context.Context, options ListTeamOptions) ([]Team, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var teamsResult teamSearchResult

	resp, err := i.Client.Do(req, &teamsResult)
	if err != nil {
		return nil, resp, err
	}
	return teamsResult.Embedded.Teams, resp, err
}

// GetTeam returns a Team with the provided teamId. Setting detailed to true will include certain hidden fields.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_team_info
func (i *IdentityService) GetTeam(ctx context.Context, teamId string) (*Team, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams/"+teamId, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var getTeam Team

	resp, err := i.Client.Do(req, &getTeam)
	if err != nil {
		return nil, resp, err
	}
	return &getTeam, resp, err
}

// CreateTeam creates a new team using the provided Team object.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_create_team
func (i *IdentityService) CreateTeam(ctx context.Context, team *Team) (*Team, *Response, error) {
	buf, err := json.Marshal(team)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams", http.MethodPost, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	var newTeam Team
	resp, err := i.Client.Do(req, &newTeam)
	if err != nil {
		return nil, resp, err
	}

	return &newTeam, resp, nil
}

// DeleteTeam deletes a team from the Veracode API using the provided teamId.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_delete_team
func (i *IdentityService) DeleteTeam(ctx context.Context, teamId string) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams/"+teamId, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// UpdateTeam updates a specific team and sets nulls to fields not in the request (if the database allows it) unless partial is set to true.
// If incremental is set to true, any values in the users list will be added to the teams's users instead of replacing them.
//
// Veracode API documentation: https://docs.veracode.com/r/c_identity_update_team
func (i *IdentityService) UpdateTeam(ctx context.Context, team *Team, options UpdateOptions) (*Team, *Response, error) {
	buf, err := json.Marshal(team)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams/"+team.TeamId, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var updatedTeam Team
	resp, err := i.Client.Do(req, &updatedTeam)
	if err != nil {
		return nil, resp, err
	}

	return &updatedTeam, resp, nil
}

// SelfListTeams returns a list of teams that the current user is a part of.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fteams/getTeams_1
func (i *IdentityService) SelfListTeams(ctx context.Context, options ListTeamOptions) ([]Team, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/teams/self", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = QueryEncode(options)

	var teamsResult teamSearchResult

	resp, err := i.Client.Do(req, &teamsResult)
	if err != nil {
		return nil, resp, err
	}
	return teamsResult.Embedded.Teams, resp, err
}
