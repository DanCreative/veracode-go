package veracode

import (
	"fmt"
	"strings"
)

// type Region string

// const (
// 	RegionEurope       Region = "https://api.veracode.eu"
// 	RegionUnitedStates Region = "https://api.veracode.us"
// 	RegionCommercial   Region = "https://api.veracode.com"
// )

var regions = map[string]string{
	"e": "https://api.veracode.eu",
	"f": "https://api.veracode.us",
	"g": "https://api.veracode.com",
}

func getRegionFromCredentials(apiKey string) (string, error) {
	var regionCharacter string
	if strings.Contains(apiKey, "-") {
		prefix := strings.Split(apiKey, "-")[0]
		if len(prefix) != 8 {
			return "", fmt.Errorf("credential %s starts with an invalid prefix", apiKey)
		}

		regionCharacter = strings.ToLower(string(prefix[6]))
	} else {
		regionCharacter = "g"
	}

	if v, ok := regions[regionCharacter]; ok {
		return v, nil
	} else {
		return "", fmt.Errorf("credential %s does not map to a known region", apiKey)
	}
}
