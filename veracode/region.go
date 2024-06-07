package veracode

type Region string

const (
	RegionEurope       Region = "https://api.veracode.eu"
	RegionUnitedStates Region = "https://api.veracode.us"
	RegionCommercial   Region = "https://api.veracode.com"
)
