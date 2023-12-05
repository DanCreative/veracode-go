package veracode

import (
	"bytes"
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
	Users        *[]User          `json:"users,omitempty"`
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

// GetTeam returns a Team with the provided teamId. Setting detailed to true will include certain hidden fields.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_team_info
func (i *IdentityService) GetTeam(ctx context.Context, teamId string) (*Team, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/teams/"+teamId, http.MethodGet, nil)
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
func (i *IdentityService) CreateTeam(ctx context.Context, team *Team) (*Team, *http.Response, error) {
	buf, err := json.Marshal(team)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/teams", http.MethodPost, bytes.NewBuffer(buf))
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
func (i *IdentityService) DeleteTeam(ctx context.Context, teamId string) (*http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/teams/"+teamId, http.MethodDelete, nil)
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
func (i *IdentityService) UpdateTeam(ctx context.Context, team *Team, options UpdateOptions) (*Team, *http.Response, error) {
	buf, err := json.Marshal(team)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/teams/"+team.TeamId, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = values.Encode()

	var updatedTeam Team
	resp, err := i.Client.Do(req, &updatedTeam)
	if err != nil {
		return nil, resp, err
	}

	return &updatedTeam, resp, nil
}
