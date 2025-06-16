package veracode

import (
	"context"
	"net/http"
	"time"
)

// ctime is used to unmarshal the different time formats returned by the API into [time.Time]
type ctime struct {
	time.Time
}

type APICredentials struct {
	ApiId          string   `json:"api_id"`
	ApiSecret      string   `json:"api_secret"`
	ExpirationTs   ctime    `json:"expiration_ts"`
	RevocationUser string   `json:"revocation_user"`
	RevocationTs   ctime    `json:"revocation_ts"`
	Links          NavLinks `json:"_links"`
}

func (ct *ctime) UnmarshalJSON(b []byte) (err error) {
	// Standard RFC 3339
	err = ct.Time.UnmarshalJSON(b)
	if err == nil {
		return nil
	}

	// Variation of ISO 8601
	n, err := time.Parse("2006-01-02T15:04:05.000Z0700", string(b[1:len(b)-1]))
	if err == nil {
		ct.Time = n
		return nil
	}

	// ISO 8601
	n, err = time.Parse("2006-01-02 15:04:05 UTC", string(b[1:len(b)-1]))
	if err == nil {
		ct.Time = n
		return nil
	}

	return err
}

// SelfGetCredentials returns the current user's API credentials.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/getUserApiCreds
func (i *IdentityService) SelfGetCredentials(ctx context.Context) (APICredentials, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials", http.MethodGet, nil)
	if err != nil {
		return APICredentials{}, nil, err
	}

	var selfCred APICredentials

	resp, err := i.Client.Do(req, &selfCred)
	if err != nil {
		return APICredentials{}, resp, err
	}
	return selfCred, resp, err
}

// SelfGenerateCredentials generates a new API credentials for the current user.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/generateApiCreds
func (i *IdentityService) SelfGenerateCredentials(ctx context.Context) (APICredentials, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials", http.MethodPost, nil)
	if err != nil {
		return APICredentials{}, nil, err
	}

	var selfCred APICredentials

	resp, err := i.Client.Do(req, &selfCred)
	if err != nil {
		return APICredentials{}, resp, err
	}

	return selfCred, resp, nil
}

// SelfRevokeCredentials revokes the current user's API credentials.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/revokeUserApiCreds
func (i *IdentityService) SelfRevokeCredentials(ctx context.Context) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials", http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// GetCredentialsByUserId returns the API credentials for the provided userId.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/getApiCredsByUserId
func (i *IdentityService) GetCredentialsByUserId(ctx context.Context, userId string) (APICredentials, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials/user_id/"+userId, http.MethodGet, nil)
	if err != nil {
		return APICredentials{}, nil, err
	}

	var selfCred APICredentials

	resp, err := i.Client.Do(req, &selfCred)
	if err != nil {
		return APICredentials{}, resp, err
	}
	return selfCred, resp, err
}

// GenerateCredentialsByUserId generates new API credentials for the provided userId.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/generateApiCredsByUserId
func (i *IdentityService) GenerateCredentialsByUserId(ctx context.Context, userId string) (APICredentials, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials/user_id/"+userId, http.MethodPost, nil)
	if err != nil {
		return APICredentials{}, nil, err
	}

	var selfCred APICredentials

	resp, err := i.Client.Do(req, &selfCred)
	if err != nil {
		return APICredentials{}, resp, err
	}

	return selfCred, resp, nil
}

// RevokeCredentialsByUserId revokes the API credentials for the provided userId.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/revokeApiCredsByUserId
func (i *IdentityService) RevokeCredentialsByUserId(ctx context.Context, userId string) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials/user_id/"+userId, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// GetCredentialsByKey returns the API credentials for the provided API key.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/getApiCreds
func (i *IdentityService) GetCredentialsByKey(ctx context.Context, Apikey string) (APICredentials, *Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials/"+Apikey, http.MethodGet, nil)
	if err != nil {
		return APICredentials{}, nil, err
	}

	var selfCred APICredentials

	resp, err := i.Client.Do(req, &selfCred)
	if err != nil {
		return APICredentials{}, resp, err
	}
	return selfCred, resp, err
}

// RevokeCredentialsByKey revokes the API credentials for the provided API key.
//
// Veracode API documentation:
//   - https://app.swaggerhub.com/apis/Veracode/veracode-identity_api/1.1#/%2Fv2%2Fapi_credentials/revokeApiCreds
func (i *IdentityService) RevokeCredentialsByKey(ctx context.Context, Apikey string) (*Response, error) {
	req, err := i.Client.NewRequest(ctx, "/api/authn/v2/api_credentials/"+Apikey, http.MethodDelete, nil)
	if err != nil {
		return nil, err
	}

	resp, err := i.Client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, err
}
