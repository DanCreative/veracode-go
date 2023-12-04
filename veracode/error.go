package veracode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// The Veracode APIs can return 2 different error response bodies.
//
// Example of a general error:
//
//	{
//	   "http_code": 400,
//	   "http_status": "Bad Request",
//	   "message": "Invalid UUID string: abcd"
//	}
//
// Example of an error with the provided query values:
//
//	{
//	  "errors": [
//	    "team_id: Invalid value.",
//	    "role_id: not a valid GUID"
//	  ],
//	  "status": 400
//	}
type Error struct {
	HttpCode   int    `json:"http_code"`
	HttpStatus string `json:"http_status"`
	Message    string `json:"message"`

	Errors []string `json:"errors"`
	Status int      `json:"status"`
}

func (v *Error) Error() string {
	if len(v.Errors) == 0 {
		return fmt.Sprintf("status: %d; message: %s", v.HttpCode, v.Message)
	} else {
		return fmt.Sprintf("status: %d; messages: [\"%s\"]", v.Status, strings.Join(v.Errors, "\", \""))
	}
}

// NewVeracodeError unmarshalls a response body into a new Veracode error.
func NewVeracodeError(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("failed to create new Veracode Error because response is empty")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	verr := &Error{}
	err = json.Unmarshal(body, verr)
	if err != nil {
		return err
	}
	return verr
}
