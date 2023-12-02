package veracode

import "encoding/json"

type Team struct {
	TeamId       string           `json:"team_id,omitempty"`
	TeamLegacyId int              `json:"team_legacy_id,omitempty"`
	TeamName     string           `json:"team_name,omitempty"`
	Relationship TeamRelationship `json:"relationship,omitempty"`
}

type TeamRelationship struct {
	Name string `json:"name,omitempty"`
}

func (t TeamRelationship) MarshalJSON() ([]byte, error) {
	jsonValue, err := json.Marshal(t.Name)
	return jsonValue, err
}
