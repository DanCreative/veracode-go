package veracode

type Region int

const (
	RegionEurope = iota
	RegionUnitedStates
	RegionCommercial
)

// parseRegion takes a region and returns the base URL for that region
func parseRegion(region Region) string {
	switch region {
	case RegionEurope:
		return "https://api.veracode.eu/api/authn/v2"
	case RegionCommercial:
		return "https://api.veracode.com/api/authn/v2"
	case RegionUnitedStates:
		return "https://api.veracode.us/api/authn/v2"
	default:
		return "https://api.veracode.com/api/authn/v2"
	}
}
