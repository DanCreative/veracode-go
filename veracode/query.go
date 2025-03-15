package veracode

import (
	"net/url"
	"reflect"
	"strings"

	"github.com/gorilla/schema"
)

var enc = newEncoder()

// QueryEncode takes any object and encodes it to a query string, while replacing "+" with "%20".
//
// The reason I added this function, was because the Veracode APIs does not support "+" to indicate spaces in the URL's query parameters.
// Example: `?name=foo+bar` will cause a 401 error.
//
// # Known bug:
//
// if "+" is part of the query parameter name/value before encoding, it will also be replaced by "%20".
// I am doing it this way for simplicity, performance (the alternative is to loop through the url.Values map and replace specifically every space before encoding) and
// because I don't currently have a use case to pass any values that contain "+".
func QueryEncode(options any) string {
	q := url.Values{}

	if enc.Encode(options, q) != nil {
		return q.Encode()
	}

	return strings.Replace(q.Encode(), "+", "%20", -1)
}

func newEncoder() *schema.Encoder {
	r := schema.NewEncoder()
	r.SetAliasTag("url")

	// Allows for asc/desc individual fields
	r.RegisterEncoder(SortQueryField{}, func(v reflect.Value) string {
		s := v.FieldByName("Name").String()

		if v.FieldByName("IsDesc").Bool() {
			s += ",desc"
		}

		return s
	})
	return r
}
