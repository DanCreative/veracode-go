package veracode

// PageMeta contains the meta data for the current API page.
type PageMeta struct {
	Number        int `json:"number"`
	Size          int `json:"size"`
	TotalElements int `json:"total_elements"`
	TotalPages    int `json:"total_pages"`
}

// Container of navigation links.
type NavLinks struct {
	First link `json:"first"`
	Last  link `json:"last"`
	Next  link `json:"next"`
	Prev  link `json:"prev"`
	Self  link `json:"self"`
}

// Link contains the URL to a milestone page.
type link struct {
	HrefURL string `json:"href"`
}

type SortQueryField struct {
	Name   string
	IsDesc bool
}

// PageOptions contains fields used to page through an endpoint as well as set page size.
type PageOptions struct {
	Size int              `url:"size,omitempty"` // Increase the page size.
	Page int              `url:"page"`           // Page through the list.
	Sort []SortQueryField `url:"sort,omitempty"` // Sort by multiple field names. Field names have to be in camelCase. Sort is ascending by default.
}
