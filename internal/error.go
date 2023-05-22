package internal

import (
	"encoding/json"
	"fmt"
)

type ApplicationError struct {
	HTTPStatusCode int    `json:"statusCode,omitempty"`
	ErrorCode      string `json:"code,omitempty"`
	ErrorText      string `json:"message,omitempty"`
}

func NewApplicationError(httpStatusCode int, errorCode, errorText string) ApplicationError {
	return ApplicationError{HTTPStatusCode: httpStatusCode, ErrorCode: errorCode, ErrorText: errorText}
}

func ApplicationErrorFromJson(body []byte, status int) ApplicationError {
	if status >= 500 {
		return ApplicationErrorInternalServerError.CustomMessage(string(body))
	}
	appErr := ApplicationError{}
	err := json.Unmarshal(body, &appErr)
	if err != nil {
		return ApplicationErrorInternalServerError.CustomMessage(fmt.Sprintf("body %s; err %s", string(body), err))
	}
	appErr.HTTPStatusCode = status
	return appErr
}

func (err ApplicationError) Error() string {
	return fmt.Sprintf("status: %d, code: %s, message: %s", err.HTTPStatusCode, err.ErrorCode, err.ErrorText)
}

func (err ApplicationError) CustomMessage(message string) ApplicationError {
	err.ErrorText = message
	return err
}

var (
	ApplicationErrorForbidden           = ApplicationError{HTTPStatusCode: 403, ErrorCode: "FORBIDDEN", ErrorText: "forbidden"}
	ApplicationErrorInvalidInput        = ApplicationError{HTTPStatusCode: 400, ErrorCode: "INVALID_INPUT", ErrorText: "invalid input"}
	ApplicationErrorInternalServerError = ApplicationError{HTTPStatusCode: 500, ErrorCode: "INTERNAL_SERVER_ERROR", ErrorText: "unknown error"}
)
