package veracode

import (
	"encoding/json"
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
