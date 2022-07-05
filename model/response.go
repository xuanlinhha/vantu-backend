package model

import (
	"errors"
	"net/http"
)

type ResponseData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

const (
	OK_STATUS string = "OK"
)

// error types
var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrNotFound            = errors.New("Your requested Item is not found")
	ErrBadParamInput       = errors.New("Given Param is not valid")
)

// http code
var httpStatusMapping = map[error]int{
	ErrNotFound:            http.StatusNotFound,
	ErrInternalServerError: http.StatusInternalServerError,
	ErrBadParamInput:       http.StatusBadRequest,
}

// app code
var appCodeMapping = map[error]string{
	ErrNotFound:            "NotFound",
	ErrInternalServerError: "InternalServerError",
	ErrBadParamInput:       "BadParamInput",
}

// convert error to http code
func GetHttpStatus(err error) int {
	if val, ok := httpStatusMapping[err]; ok {
		return val
	}
	return http.StatusInternalServerError
}

// convert error to app code
func GetAppCode(err error) string {
	if val, ok := appCodeMapping[err]; ok {
		return val
	}
	return OK_STATUS
}
