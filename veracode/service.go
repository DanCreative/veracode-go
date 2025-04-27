package veracode

// service is the base structure to bundle API services under a sub-struct.
type service struct {
	Client *Client
}

// You can use the Identity Service to manage the administrative configuration for your organization that is in the Veracode Platform.
// For more information: https://docs.veracode.com/r/c_identity_intro.
//
// Currently supports V2 of the Identity API
type IdentityService service

// You can use the Applications API to quickly access information about your Veracode applications.
// For more information, review the documentation: https://docs.veracode.com/r/c_apps_intro
//
// Currently supports V1 of the Applications API
type ApplicationService service

// You can use the Development Sandbox API to create, update, and delete development sandboxes.
// For more information:
// 	- https://docs.veracode.com/r/c_rest_sandbox_intro
// 	- https://app.swaggerhub.com/apis/Veracode/veracode-development_sandbox_api/2.0#/Application%20Sandbox%20Information%20API
//
// Currently supports V1 of the Development Sandbox API
type SandboxService service

// You can use the Healthcheck API to perform a simple test for verifying authenticated connectivity to Veracode.
//
// The Healthcheck API provides this lightweight endpoint: /healthcheck/status
//
// You use the endpoint to verify that Veracode services are available and responding to authentication events, instead of using other API calls that can potentially return large volumes of data.
//
// For more information: https://docs.veracode.com/r/c_healthcheck_intro
type HealthCheckService service
