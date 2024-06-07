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

	rateTransport, err := veracode.NewRateTransport(nil, time.Second, 10)
	check(err)

	jar, err := cookiejar.New(&cookiejar.Options{})
	check(err)

	httpClient := &http.Client{
		Transport: rateTransport,
		Jar:       jar,
	}

	client, err := veracode.NewClient(veracode.RegionEurope, httpClient, key, secret)
	check(err)

	ctx := context.Background()

	teams, resp, err := client.Identity.ListTeams(ctx, veracode.ListTeamOptions{Size: 10})
	check(err)

	fmt.Println(teams)
	fmt.Println(resp)
}
```

## Release Notes:
### Version ```0.2.0```:
#### General:
- ```Region``` is now just a type definition of ```String```. This change allows new regions to be added without requiring the package to be updated.
- Added functionality to update the region hostname in a concurrency-safe way.
#### Identity API v2:
- Added a new ```RoleUser``` struct to represent the roles as part of the ```User``` aggregate struct. This change makes it more clear which role fields are available when calling different endpoints.
### Version ```0.1.0```:
#### General:
- Added functionality to load credentials from the credentials file and swap between profiles.
- HMAC is handled using my [veracode-hmac-go](https://github.com/DanCreative/veracode-hmac-go) package.
- Calling code can add additional Transports to the HTTP client. In above example, a rate limiter is added. When the Client is created, it automatically daisy-chains the authentication Transport to the provided Transport(s).
- The client exposes several functions to allow the calling code to implement any endpoints not already available. Namely: ```NewRequest()``` and ```Do()```.
-  All of the page meta data for collection requests are returned in the ```Response``` struct, which wraps the ```http.Response``` struct.
#### Identity API v2:
- Added support for user, team, business-unit and role endpoints.