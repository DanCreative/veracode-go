package veracode

// service is the base structure to bundle API services under a sub-struct.
type service struct {
	Client *Client
}

// You can use the Identity Service to manage the administrative configuration for your organization that is in the Veracode Platform.
// For more information: https://docs.veracode.com/r/c_identity_intro.
type IdentityService service
