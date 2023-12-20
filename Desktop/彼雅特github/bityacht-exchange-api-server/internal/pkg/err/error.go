package errpkg

import (
	"encoding/json"
	"fmt"
)

// Error is Common error type for this project.
type Error struct {
	HttpStatus int
	Code       Code
	Data       interface{}
	Err        error
}

// JsonError for swagger document and JSON Marshaler
type JsonError struct {
	Code  Code        `json:"code" binding:"required"`
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func (e Error) MarshalJSON() ([]byte, error) {
	if e.Err == nil {
		return json.Marshal(JsonError{Code: e.Code, Data: e.Data})
	}

	switch err := e.Err.(type) {
	case *Error:
		if err == nil {
			return json.Marshal(JsonError{Code: e.Code, Data: e.Data})
		}
	}

	return json.Marshal(JsonError{Code: e.Code, Data: e.Data, Error: e.Err.Error()})
}

func (e Error) Error() string {
	var output string

	if e.Code == 0 {
		output = fmt.Sprintf("code: %+v", CodeUnknown)
	} else {
		output = fmt.Sprintf("code: %+v", e.Code)
	}

	if e.Data != nil {
		output += fmt.Sprintf(", data: %+v", e.Data)
	}
	if e.Err != nil {
		switch err := e.Err.(type) {
		case *Error:
			if err != nil && err.Err != nil {
				output += ", err: " + err.Err.Error()
			}
		default:
			output += ", err: " + e.Err.Error()
		}
	}

	return output
}

func (e *Error) CodeEqualTo(code Code) bool {
	if e == nil {
		return code == 0
	}

	return e.Code == code
}
