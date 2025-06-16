package veracode

import (
	"fmt"
	"strings"
)

type Region map[string]string

var Regions = map[string]map[string]string{
	"e": Region{"rest": "https://api.veracode.eu", "xml": "https://analysiscenter.veracode.eu"},
	"f": Region{"rest": "https://api.veracode.us", "xml": "https://analysiscenter.veracode.us"},
	"g": Region{"rest": "https://api.veracode.com", "xml": "https://analysiscenter.veracode.com"},
}

func GetRegionFromCredentials(apiKey string) (Region, error) {
	var regionCharacter string
	if strings.Contains(apiKey, "-") {
		prefix := strings.Split(apiKey, "-")[0]
		if len(prefix) != 8 {
			return nil, fmt.Errorf("credential %s starts with an invalid prefix", apiKey)
		}

		regionCharacter = strings.ToLower(string(prefix[6]))
	} else {
		regionCharacter = "g"
	}

	if v, ok := Regions[regionCharacter]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("credential %s does not map to a known region", apiKey)
	}
}
