package veracode

import (
	"reflect"
	"testing"
)

func TestUser_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		u       *User
		want    []byte
		wantErr bool
	}{
		{
			name: "relationship with name",
			u: &User{
				UserName: "Test User",
				Relationship: TeamRelationship{
					Name: "MEMBER",
				},
			},
			want:    []byte(`{"user_name":"Test User","relationship":"MEMBER"}`),
			wantErr: false,
		},
		{
			name: "relationship with no name",
			u: &User{
				UserName:     "Test User",
				Relationship: TeamRelationship{},
			},
			want:    []byte(`{"user_name":"Test User"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.u.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("User.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}
