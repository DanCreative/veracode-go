package veracode

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
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
	Detailed     string   `url:"detailed,omitempty"`              // Passing detailed will return additional hidden fields. Value should be one of: Yes or No
	Page         int      `url:"page,omitempty"`                  // Page through the list.
	Size         int      `url:"size,omitempty"`                  // Increase the page size.
	UserName     string   `url:"user_name,omitempty"`             // Filter by username. You must specify the full username. The request does not support matching partial usernames.
	EmailAddress []string `url:"email_address,omitempty" del:","` // Filter by email address(es).
}

type SearchUserOptions struct {
	Detailed     string `url:"detailed,omitempty"`      // Passing detailed will return additional hidden fields. Value should be one of: Yes or No
	Page         int    `url:"page,omitempty"`          // Page through the list.
	Size         int    `url:"size,omitempty"`          // Increase the page size.
	SearchTerm   string `url:"search_term,omitempty"`   // You can search for partial strings of the username, first name, last name, or email address.
	RoleId       string `url:"role_id,omitempty"`       // Filter users by their role. Value should be a valid Role Id.
	UserType     string `url:"user_type,omitempty"`     // Filter by user type. Value should be one of: user or api
	LoginEnabled string `url:"login_enabled,omitempty"` // Filter by whether the login is enabled. Value should be one of: Yes or No
	LoginStatus  string `url:"login_status,omitempty"`  // Filter by the login status. Value should be one of: Active, Locked or Never
	SamlUser     string `url:"saml_user,omitempty"`     // Filter by whether the user is a SAML user or not. Value should be one of: Yes or No
	TeamId       string `url:"team_id,omitempty"`       // Filter users by team membership. Value should be a valid Team Id.
	ApiId        string `url:"api_id,omitempty"`        // Filter user by their API Id.
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

// GetUser returns user with provided userId. Setting detailed to true will include certain hidden fields.
func (i *IdentityService) GetUser(ctx context.Context, userId string, detailed bool) (User, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/users/"+userId, http.MethodGet, nil)
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

// ListUsers takes a ListUserOptions and returns a list of users. Veracode API documentation: https://docs.veracode.com/r/c_identity_list_users.
func (i *IdentityService) ListUsers(ctx context.Context, options ListUserOptions) ([]User, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/users", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = values.Encode()

	var usersResult userSearchResult

	resp, err := i.Client.Do(req, &usersResult)
	if err != nil {
		return nil, resp, err
	}
	return usersResult.Embedded.Users, resp, err
}

// SearchUsers takes a SearchUserOptions and returns a list of users. Veracode API documentation: https://docs.veracode.com/r/c_identity_search_users.
func (i *IdentityService) SearchUsers(ctx context.Context, options SearchUserOptions) ([]User, *http.Response, error) {
	req, err := i.Client.NewRequest(ctx, "/users/search", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}

	req.URL.RawQuery = values.Encode()

	var usersResult userSearchResult

	resp, err := i.Client.Do(req, &usersResult)
	if err != nil {
		return nil, resp, err
	}
	return usersResult.Embedded.Users, resp, err
}

// UpdateUser updates an existing Veracode user. Passing true to the isPartial parameter, will only update
// the fields on the user that are present in the request body. API documentation: https://docs.veracode.com/r/c_identity_update_user.
func (i *IdentityService) UpdateUser(ctx context.Context, user User, isPartial bool) (*http.Response, error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	req, err := i.Client.NewRequest(ctx, "/users/"+user.UserId, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, err
	}

	if isPartial {
		req.URL.RawQuery = "partial=true"
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}
