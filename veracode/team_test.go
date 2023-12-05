package veracode

import (
	"reflect"
	"testing"
)

func TestTeam_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		tr      *Team
		want    []byte
		wantErr bool
	}{
		{
			name: "relationship with name",
			tr: &Team{
				TeamName: "Test Team",
				Relationship: TeamRelationship{
					Name: "MEMBER",
				},
			},
			want:    []byte(`{"team_name":"Test Team","relationship":"MEMBER"}`),
			wantErr: false,
		},
		{
			name: "relationship with no name",
			tr: &Team{
				TeamName:     "Test Team",
				Relationship: TeamRelationship{},
			},
			want:    []byte(`{"team_name":"Test Team"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.tr.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Team.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Team.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}
