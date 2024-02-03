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

## Features
- Authentication using your Veracode credentials file. HMAC is handled using my [veracode-hmac-go](https://github.com/DanCreative/veracode-hmac-go) project.
- Library currently supports most of the endpoints in the Veracode Identity API.