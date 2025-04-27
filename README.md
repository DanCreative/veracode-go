## About
Go API client for the Veracode platform.

This library is still in early stages of development. It will be updated as I require features in other Veracode related projects.

## Installation
```
go get -u github.com/DanCreative/veracode-go
```

## Example
```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/DanCreative/veracode-go/veracode"
)

func main() {
	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	key, secret, err := veracode.LoadVeracodeCredentials()
	check(err)

	jar, err := cookiejar.New(&cookiejar.Options{})
	check(err)

	httpClient := &http.Client{
		Jar: jar,
	}

	client, err := veracode.NewClient(httpClient, key, secret)
	check(err)

	ctx := context.Background()

	teams, resp, err := client.Identity.ListTeams(ctx, veracode.ListTeamOptions{Size: 10})
	check(err)

	fmt.Println(teams)
	fmt.Println(resp)
}
```

## Implementation Status

> [!NOTE]
> **Legend:**
> 
> <table><tr><td>🟢</td><th>Implemented</th></tr><tr><td>🟠</td><th>Partially Implemented</th></tr><tr><td>🔴</td><th>Not Implemented Yet</th></tr><tr><td>🔷</td><th>Planned Next</th></tr><tr><td>🔻</td><th>Not Planned / In Scope</th></tr><tr></table>
>

| Service | Version | Entity | Status | Priority |
| --- | --- | --- | --- | --- |
| Identity | `v2` `latest` | user | 🟢 | |
|  |  | business unit | 🟢 | |
|  |  | api credentials | 🟢 | |
|  |  | team | 🟢 | |
|  |  | role | 🟠 | |
|  |  | jit default settings | 🔴 | 🔻 |
|  |  | permissions | 🔴 | |
| Application | `v1` `latest` | application | 🟢 | |
|  | | collection | 🟢 | |
| | | custom fields | 🟠 | |
| Development Sandbox | `v1` `latest` | sandbox | 🟢 | |
| Healthcheck | `na` | | 🟢 | |
| Policy | `v1` `latest`  | policy (incl. version) | 🔴 | 🔷 |
|  |   | policy settings | 🔴 | 🔷 |
|  |   | policy license | 🔴 | 🔷 |
| Annotations | `v2` `latest` |  | 🔴 |  |
| DAST | `v1` `latest` |  | 🔴 |  |
| DAST Essentials | `v1` `latest` |  | 🔴 |  |
| eLearning | `v1` `latest` |  | 🔴 | 🔻 |
| Findings | `v2` `latest` |  | 🔴 | |
| Greenlight | `v3` `latest` |  | 🔴 | 🔻 |
| Pipeline | `v1` `latest` |  | 🔴 | 🔻 |
| Reporting | `v1` `latest` |  | 🔴 | |
| SCA | `v1-v3` `latest` |  | 🔴 | 🔻 |
| Security Labs |  `na` |  | 🔴 | 🔻 |
| Summary Report | `v2` `latest` | summary report | 🔴 | |

## Custom Endpoints

If the endpoint that you need to call is not currently implemented, you can implement it yourself using the Client's helper function. To do so, you can wrap the Client into a custom local Client struct. Please see an example below:

```go
// Entity in this example, is the model that you will be requesting.
type Entity struct {
	Name string
}

// EntityOptions in this example, is the list options. These options will be marshalled into the query parameters.
type EntityOptions struct {
	Size int `url:"size,omitempty"`
	Page int `url:"page"`
}

// entitySearchResult in this example, is the model that will contain all of the entities in the list.
// For collection result models, make sure that the struct implements below interface:
/*
	type CollectionResult interface {
		GetLinks() navLinks
		GetPageMeta() pageMeta
	}
*/
//
// That will allow the Client to retrieve the meta data and add it to the returned veracode.Response struct.
type entitySearchResult struct {
	Embedded struct {
		Entities []Entity `json:"entities"`
	} `json:"_embedded"`
	Links veracode.NavLinks `json:"_links"`
	Page  veracode.PageMeta `json:"page"`
}

func (r *entitySearchResult) GetLinks() veracode.NavLinks {
	return r.Links
}

func (r *entitySearchResult) GetPageMeta() veracode.PageMeta {
	return r.Page
}

// Client wraps the veracode.Client.
type Client struct {
	*veracode.Client
}

// Example of requesting a single entity.
func (c *Client) GetEntity(ctx context.Context, entityGuid string) (*Entity, *veracode.Response, error) {
	// veracode.Client.NewRequest() is a helper method that creates a new request with the full resolved
	// absolute path of the provided endpoint path.
	req, err := c.NewRequest(ctx, fmt.Sprintf("/path/to/entities/%s", entityGuid), http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	var result Entity

	// veracode.Client.Do() is a helper method that executes the provided http.Request, handles the authentication and marshals the 
	// JSON response body into either the provided struct or into an error if an error occurred.
	resp, err := c.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return &result, resp, nil
}

// Example of requesting a list of entities.
func (c *Client) ListEntity(ctx context.Context, options EntityOptions) ([]Entity, *veracode.Response, error) {
	req, err := c.NewRequest(ctx, "/path/to/entities", http.MethodGet, nil)
	if err != nil {
		return nil, nil, err
	}

	// helper function veracode.QueryEncode() encodes the options into query parameters.
	// It also handles some Veracode API specific behaviours.
	req.URL.RawQuery = veracode.QueryEncode(options)

	var result entitySearchResult

	resp, err := c.Do(req, &result)
	if err != nil {
		return nil, resp, err
	}

	return result.Embedded.Entities, resp, nil
}
```

## Release Notes:

### Version ```0.6.0```
- Added Healthcheck endpoint.
- Added Development Sandboxes endpoints.
- Minor QOL updates.

### Version ```0.5.1```
<details>
<summary>See Details</summary>
- Added function to automatically determine the region based on the API credentials. Function is based on code from the python veracode-api-signing package.
- Added method to the ```Client``` to change API credentials after initialization.
- Added functions for all API credential endpoints.
- Added sorting option for the Identity service endpoints.
- Added function to get all users not in team.
- Added function to get all teams that the current user is a part of.
- Bug fixes.
</details>

### Version ```0.4.0```
<details>
<summary>See Details</summary>

#### General:
- Moved Module https://github.com/DanCreative/veracode-hmac-go into this module as a package (finally).
- Added a LICENSE file to the repository. This project is going to be using the MIT license.
- Merged the rate limiting and authentication transports into a single struct and added a default implementation.
- All collection-of-entity structs now need to implement the CollectionResult interface in order to get the navigational links and page meta details:
	```go
	type CollectionResult interface {
		GetLinks() navLinks
		GetPageMeta() pageMeta
	}
	```
	This resolves a previous issue where all collection structs needed to be added to a switch in order to get this information.
- Added support for unmarshalling all of the different error models that can be returned by the APIs.
- Fixed an issue with the Veracode API not supporting "+" as an encoding for spaces in the query string. See the veracode/query.go file for more information.

#### Application API v1:
- Added CRUD support for Applications.
- Added CRUD support for Collections.
- Added function to get a list of the custom fields.
</details>

### Version ```0.3.0```:
<details>
<summary>See Details</summary>

#### General:
- Added functionality to get different profiles from the credentials file.
</details>


### Version ```0.2.0```:
<details>
<summary>See Details</summary>

#### General:
- ```Region``` is now just a type definition of ```String```. This change allows new regions to be added without requiring the package to be updated.
- Added functionality to update the region hostname in a concurrency-safe way.
#### Identity API v2:
- Added a new ```RoleUser``` struct to represent the roles as part of the ```User``` aggregate struct. This change makes it more clear which role fields are available when calling different endpoints.

</details>




### Version ```0.1.0```:
<details>
<summary>See Details</summary>

#### General:
- Added functionality to load credentials from the credentials file and swap between profiles.
- HMAC is handled using my [veracode-hmac-go](https://github.com/DanCreative/veracode-hmac-go) package.
- Calling code can add additional Transports to the HTTP client. In above example, a rate limiter is added. When the Client is created, it automatically daisy-chains the authentication Transport to the provided Transport(s).
- The client exposes several functions to allow the calling code to implement any endpoints not already available. Namely: ```NewRequest()``` and ```Do()```.
-  All of the page meta data for collection requests are returned in the ```Response``` struct, which wraps the ```http.Response``` struct.
#### Identity API v2:
- Added support for user, team, business-unit and role endpoints.

</details>