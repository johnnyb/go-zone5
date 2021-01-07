package zone5

import (
	"net/http"
)

type ErrorWithData struct {
	Details map[string]interface{}
	HTTPStatusCode int
	Code string
	Message string
}

func (e *ErrorWithData) Error() string {
	return e.Message
}

func NewErrorWithData(code, message string, details map[string]interface{}) (*ErrorWithData) {
	err := ErrorWithData{
		Details: details,
		Code: code,
		Message: message,
	}

	return &err
}

func NewErrorWithResponse(resp *http.Response) (*ErrorWithData) {
	data, otherErr := unmarshalJsonFromReader(resp.Body)
	err := ErrorWithData {
		HTTPStatusCode: resp.StatusCode,
		Details: data,
	}
	if otherErr != nil {
		err.Message = "Couldn't read error message: " + otherErr.Error()
		return &err
	}

	msg := data["message"].(string)
	if msg != "" {
		err.Message = data["message"].(string)
	} else {
		err.Message = "Unknown error"
	}

	return &err
}
