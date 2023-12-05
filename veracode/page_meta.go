package veracode

// PageMeta contains the meta data for the current API page.
type pageMeta struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

// Container of navigation links.
type navLinks struct {
	First link `json:"first"`
	Last  link `json:"last"`
	Next  link `json:"next"`
	Self  link `json:"self"`
}

// Link contains the URL to a milestone page.
type link struct {
	HrefURL string `json:"href_url"`
}

// PageOptions contains fields used to page through an endpoint as well as set page size.
type PageOptions struct {
	Size int `url:"size,omitempty"`
	Page int `url:"page"`
}
