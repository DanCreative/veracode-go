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
		return "https://api.veracode.eu"
	case RegionCommercial:
		return "https://api.veracode.com"
	case RegionUnitedStates:
		return "https://api.veracode.us"
	default:
		return "https://api.veracode.com"
	}
}
