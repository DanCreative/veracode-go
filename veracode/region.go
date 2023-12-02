package veracode

type Region int

const (
	Europe = iota
	UnitedStates
	Commercial
)

// parseRegion takes a region and returns the base URL for that region
func parseRegion(region Region) string {
	switch region {
	case Europe:
		return "https://api.veracode.eu/api/authn/v2"
	case Commercial:
		return "https://api.veracode.com/api/authn/v2"
	case UnitedStates:
		return "https://api.veracode.us/api/authn/v2"
	default:
		return "https://api.veracode.com/api/authn/v2"
	}
}
