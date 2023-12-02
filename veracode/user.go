package veracode

import (
	"context"
	"net/http"
)

// userSearchResult is required to decode the list user and search user response bodies.
type userSearchResult struct {
	Embedded struct {
		Users []User `json:"users"`
	} `json:"_embedded"`
	Links navLinks `json:"_links"`
	Page  pageMeta `json:"page"`
}

type User struct {
	// Below fields will be included in /users and /users/search calls
	LoginEnabled bool   `json:"login_enabled"`
	SamlUser     bool   `json:"saml_user"`
	EmailAddress string `json:"email_address,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	LegacyUserId string `json:"legacy_user_id,omitempty"`
	UserName     string `json:"user_name,omitempty"`

	// AccountType is added by passing detailed=true in the URL values.
	// AccountType will be shown in the user model for /users/{id}, /users and /users/search
	AccountType string `json:"account_type,omitempty"`

	// Below fields will only be included in /users/{id} calls
	// BACKLOG: Add remaining fields for model as required.
	Active bool   `json:"active"`
	Roles  []Role `json:"roles,omitempty"`
	Teams  []Team `json:"teams,omitempty"`
}

type ListUserOptions struct {
	// Enabled values should be { 0: omit, 1: false, 2:true}
	Enabled      int
	Page         int
	Size         int
	UserName     string
	EmailAddress []string
}

// Self returns the requesting user's details. Setting detailed to true will add certain hidden fields.
func (i *IdentityService) Self(ctx context.Context, detailed bool) (User, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/users/self", http.MethodGet, nil)
	if err != nil {
		return User{}, nil, err
	}

	if detailed {
		req.URL.RawQuery = "detailed=true"
	}

	var user User

	resp, err := i.Client.Do(req, &user)
	if err != nil {
		return User{}, resp, err
	}
	return user, resp, err
}

// ListUsers takes a ListUserOptions and returns a list of users. For additional information please see
// the documentation: https://docs.veracode.com/r/c_identity_list_users.
func (i *IdentityService) ListUsers(ctx context.Context, options ListUserOptions) ([]User, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/users", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var usersResult userSearchResult

	resp, err := i.Client.Do(req, &usersResult)
	if err != nil {
		return nil, resp, err
	}
	return usersResult.Embedded.Users, resp, err
}
