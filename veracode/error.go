package veracode

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// The Veracode APIs can return multiple different error response bodies.
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
// Example of an error returned by the Applications API:
//
// {
//     "_embedded" : {
//       "api_errors" : [ {
//         "id" : "abcd",
//         "code" : "NOT_FOUND",
//         "title" : "The requested application could not be found",
//         "status" : "404",
//         "source" : {
//           "pointer" : "/appsec/v1/applications/abcd",
//           "parameter" : null
//         }
//       } ]
//     },
//     "fallback_type" : null,
//     "full_type" : null
// }
//
// {
// 	"id" : "blah",
// 	"code" : "BAD_REQUEST",
// 	"title" : "Failed to convert value of type 'java.lang.String' to required type 'com.veracode.agora.collection.domain.CollectionComplianceStatus'; nested exception is org.springframework.core.convert.ConversionFailedException: Failed to convert from type [java.lang.String] to type [@io.swagger.v3.oas.annotations.Parameter @org.springframework.web.bind.annotation.RequestParam com.veracode.agora.collection.domain.CollectionComplianceStatus] for value 'blah'; nested exception is java.lang.IllegalArgumentException: No enum constant com.veracode.agora.collection.domain.CollectionComplianceStatus.blah",
// 	"detail" : null,
// 	"status" : "400",
// 	"source" : {
// 	  "pointer" : "/appsec/v1/collections",
// 	  "parameter" : null
// 	}
//  }

type errorSource struct {
	Parameter string `json:"parameter,omitempty"`
	Pointer   string `json:"pointer,omitempty"`
}

type errorBodyJson struct {
	HttpCode   int    `json:"http_code"`
	HttpStatus string `json:"http_status"`
	Message    string `json:"message"`

	ErrorDescription string `json:"error_description"`
	Error            string `json:"error"`

	Title  string   `json:"title,omitempty"`
	Errors []string `json:"errors"`
	// Status int      `json:"status"`

	Detail string      `json:"detail,omitempty"`
	Source errorSource `json:"source,omitempty"`

	Embedded struct {
		Errors []struct {
			Id     string `json:"id,omitempty"`
			Code   string `json:"code,omitempty"`
			Title  string `json:"title,omitempty"`
			Status string `json:"status,omitempty"`
			Source struct {
				Pointer   string `json:"pointer,omitempty"`
				Parameter string `json:"parameter,omitempty"`
			} `json:"source,omitempty"`
		} `json:"api_errors,omitempty"`
	} `json:"_embedded,omitempty"`
}

type errorBodyXml struct {
	XMLName xml.Name `xml:"error"`
	Message string   `xml:",chardata"`
}

type Error struct {
	Code     int
	Endpoint string
	Messages []string
}

func (v Error) Error() string {
	return fmt.Sprintf("api error returned from %s (%d): [\"%s\"]", v.Endpoint, v.Code, strings.Join(v.Messages, "\", \""))
}

func (e *Error) UnmarshalJSON(data []byte) (err error) {
	errBody := errorBodyJson{}
	err = json.Unmarshal(data, &errBody)
	if err != nil {
		return
	}

	if errBody.ErrorDescription != "" {
		e.Messages = append(e.Messages, fmt.Sprintf("description:'%s'", errBody.ErrorDescription), fmt.Sprintf("error:'%s'", errBody.Error))
	}

	if errBody.Detail != "" {
		e.Messages = append(e.Messages, fmt.Sprintf("detail:'%s'", errBody.Detail))
	}

	if errBody.Title != "" {
		e.Messages = append(e.Messages, fmt.Sprintf("title:'%s'", errBody.Title))
	}

	if errBody.Message != "" {
		e.Messages = append(e.Messages, fmt.Sprintf("message:'%s'", errBody.Message))
	}

	if errBody.Source.Parameter != "" {
		e.Messages = append(e.Messages, fmt.Sprintf("parameter: '%s' at pointer: '%s'", errBody.Source.Parameter, errBody.Source.Pointer))
	}

	if l := len(errBody.Errors); l > 0 {
		copy(e.Messages, errBody.Errors)
	}

	if l := len(errBody.Embedded.Errors); l > 0 {
		for _, apiError := range errBody.Embedded.Errors {
			e.Messages = append(e.Messages, fmt.Sprintf("{id:'%s', code:'%s', title: '%s', status: '%s'}", apiError.Id, apiError.Code, apiError.Title, apiError.Status))
		}
	}

	return
}

func (e *Error) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var xml errorBodyXml

	if err := d.DecodeElement(&xml, &start); err != nil {
		return err
	}

	e.Messages = []string{xml.Message}
	return nil
}

// NewVeracodeError unmarshals a response body into a new Veracode error.
func NewVeracodeError(resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("failed to create new Veracode Error because response is empty")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	verr := Error{Code: resp.StatusCode, Endpoint: resp.Request.URL.Path}
	if len(body) > 0 {
		err = json.Unmarshal(body, &verr)
		if err != nil {
			return err
		}

	}
	return verr
}
