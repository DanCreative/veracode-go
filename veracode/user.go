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
	LoginEnabled *bool  `json:"login_enabled,omitempty"`
	SamlUser     *bool  `json:"saml_user,omitempty"` // Required when creating a new SAML user.
	EmailAddress string `json:"email_address,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	UserId       string `json:"user_id,omitempty"`
	LegacyUserId string `json:"legacy_user_id,omitempty"`
	UserName     string `json:"user_name,omitempty"`

	// AccountType is added by passing detailed=true in the URL values.
	// AccountType will be shown in the user model for /users/{id}, /users and /users/search
	AccountType string `json:"account_type,omitempty"`

	Relationship TeamRelationship `json:"relationship,omitempty"` // Only present when the user is included in the Team model.

	// Below fields will only be included in /users/{id} calls
	// BACKLOG: Add remaining fields for model as required.
	Active *bool `json:"active,omitempty"`

	Roles       *[]Role       `json:"roles,omitempty"`       // Be careful when setting a user's roles to an empty list. This will remove even the Administrator role.
	Teams       *[]Team       `json:"teams,omitempty"`       // Giving a user the team admin role will require setting the Team.Relationship.Name to "ADMIN"
	Permissions *[]Permission `json:"permissions,omitempty"` // A permission with name: "apiUser" needs to be set to create a new API user.

	Title       string `json:"title,omitempty"`        // Can be set when creating a new user, but is not available when fetching a user.
	UserType    string `json:"user_type,omitempty"`    // Required when creating a new user.
	SamlSubject string `json:"saml_subject,omitempty"` // Required when creating a new SAML user.
}

type Permission struct {
	Name string `json:"permission_name,omitempty"`
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

type UpdateOptions struct {
	Incremental *bool `url:"incremental,omitempty"` // incremental=true indicates that you are adding items to a list for an object property, such as adding users to a team.
	Partial     *bool `url:"partial,omitempty"`     // partial=true indicates that you are updating only a subset of properties for an object.
}

// If Relationship.Name is "", create custom struct where TeamRelationship is a pointer and set it to nil.
// This will omit relationship from the marshalled json.
//
// If Relationship.Name is not "", flatten TeamRelationship to Relationship in User model.
func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	if u.Relationship.Name == "" {
		return json.Marshal(&struct {
			*Alias
			Relationship *TeamRelationship `json:"relationship,omitempty"`
		}{
			Alias:        (*Alias)(u),
			Relationship: nil,
		})
	}
	return json.Marshal(&struct {
		*Alias
		Relationship string `json:"relationship,omitempty"`
	}{
		Alias:        (*Alias)(u),
		Relationship: u.Relationship.Name,
	})
}

// NewUser is a helper function that creates a new user with all of the required fields to Post to the Veracode API.
//
// Note that NewUser adds the "securityinsightsonly" role as the default role for the created user. The caller should update the roles
// on the User.
func NewUser(emailAddress, firstName, lastName string) *User {
	isActive := true
	return &User{
		EmailAddress: emailAddress,
		Active:       &isActive,
		UserName:     emailAddress,
		FirstName:    firstName,
		LastName:     lastName,
		Roles:        &[]Role{{RoleName: "securityinsightsonly"}},
		UserType:     "VOSP",
	}
}

// NewSAMLUser is a helper function that creates a new SAML user with all of the required fields to Post successfully to the Veracode API.
//
// Note that NewSAMLUser adds the "securityinsightsonly" role as the default role for the created user. The caller should update the roles
// on the returned User.
func NewSAMLUser(emailAddress, firstName, lastName, samlSubject string) *User {
	isSAMLUser := true
	isActive := true
	return &User{
		EmailAddress: emailAddress,
		Active:       &isActive,
		SamlUser:     &isSAMLUser,
		SamlSubject:  samlSubject,
		FirstName:    firstName,
		LastName:     lastName,
		Roles:        &[]Role{{RoleName: "securityinsightsonly"}},
		UserType:     "VOSP",
	}
}

// NewAPIUser is a helper function that creates a new service account user with all of the required fields to Post successfully to the Veracode API.
//
// Note the following:
//   - NewAPIUser adds the "resultsapi" role as the default role for the created user. The caller should update the roles on the returned User.
//   - Providing a nil value or an empty slice for parameter "teams", will add the "noteamrestrictionapi" role for the user.
//     This role allows a Service Account to see all applications across the Veracode instance.
func NewAPIUser(userName, emailAddress, firstName, lastName string, teams []Team) *User {
	isActive := true
	user := &User{
		EmailAddress: emailAddress,
		UserName:     userName,
		Active:       &isActive,
		FirstName:    firstName,
		LastName:     lastName,
		Roles:        &[]Role{{RoleName: "resultsapi"}},
		Permissions:  &[]Permission{{Name: "apiUser"}},
		Teams:        &teams,
	}

	if len(teams) == 0 {
		newRoles := append([]Role{{RoleName: "noteamrestrictionapi"}}, *user.Roles...)
		user.Roles = &newRoles
	}
	return user
}

// Self returns the requesting user's details. Setting detailed to true will add certain hidden fields.
func (i *IdentityService) Self(ctx context.Context, detailed bool) (*User, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/self", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	if detailed {
		req.URL.RawQuery = "detailed=true"
	}

	var selfUser User

	resp, err := i.Client.Do(req, &selfUser)
	if err != nil {
		return nil, resp, err
	}
	return &selfUser, resp, err
}

// GetUser returns user with provided userId. Setting detailed to true will include certain hidden fields.
func (i *IdentityService) GetUser(ctx context.Context, userId string, detailed bool) (*User, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/"+userId, http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	if detailed {
		req.URL.RawQuery = "detailed=true"
	}

	var getUser User

	resp, err := i.Client.Do(req, &getUser)
	if err != nil {
		return nil, resp, err
	}
	return &getUser, resp, err
}

// ListUsers takes a ListUserOptions and returns a list of users.
//
// Veracode API documentation: https://docs.veracode.com/r/c_identity_list_users.
func (i *IdentityService) ListUsers(ctx context.Context, options ListUserOptions) ([]User, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users", http.MethodGet, nil)
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

// SearchUsers takes a SearchUserOptions and returns a list of users.
//
// Veracode API documentation: https://docs.veracode.com/r/c_identity_search_users.
func (i *IdentityService) SearchUsers(ctx context.Context, options SearchUserOptions) ([]User, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/search", http.MethodGet, nil)
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

// UpdateUser updates a specific user and sets nulls to fields not in the request (if the database allows it) unless partial is set to true.
// If incremental is set to true, any values in the roles or teams list will be added to the user's roles/teams instead of replacing them.
//
// Veracode API documentation: https://docs.veracode.com/r/c_identity_update_user.
func (i *IdentityService) UpdateUser(ctx context.Context, user *User, options UpdateOptions) (*User, *Response, error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/"+user.UserId, http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = values.Encode()

	var updatedUser User
	resp, err := i.Client.Do(req, &updatedUser)
	if err != nil {
		return nil, resp, err
	}

	return &updatedUser, resp, nil
}

// UpdateSelf updates the requesting user and sets nulls to fields not in the request (if the database allows it) unless partial is set to true.
// If incremental is set to true, any values in the roles or teams list will be added to the user's roles/teams instead of replacing them.
//
// Veracode API documentation: https://docs.veracode.com/r/c_identity_update_user.
func (i *IdentityService) UpdateSelf(ctx context.Context, user *User, options UpdateOptions) (*User, *Response, error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/self", http.MethodPut, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	values, err := query.Values(options)
	if err != nil {
		return nil, nil, err
	}
	req.URL.RawQuery = values.Encode()

	var updatedUser User
	resp, err := i.Client.Do(req, &updatedUser)
	if err != nil {
		return nil, resp, err
	}

	return &updatedUser, resp, nil
}

// CreateUser creates a new user using the provided User object. Setting generateApiCredentials to true, will generate API credentials for
// the user on creation.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_create_api
//   - https://docs.veracode.com/r/c_identity_create_human
func (i *IdentityService) CreateUser(ctx context.Context, user *User, generateApiCredentials bool) (*User, *Response, error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return nil, nil, err
	}

	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users", http.MethodPost, bytes.NewBuffer(buf))
	if err != nil {
		return nil, nil, err
	}

	if generateApiCredentials {
		req.URL.RawQuery = "generate_api_creds=true"
	}

	var newUser User
	resp, err := i.Client.Do(req, &newUser)
	if err != nil {
		return nil, resp, err
	}

	return &newUser, resp, nil
}

// DeleteUser deletes a user using the provided userId. This applies to both a user account and an API service account.
//
// Veracode API documentation:
//   - https://docs.veracode.com/r/c_identity_delete
func (i *IdentityService) DeleteUser(ctx context.Context, userId string) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/users/"+userId, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
