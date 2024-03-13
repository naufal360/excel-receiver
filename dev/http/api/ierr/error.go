package ierr

import (
	"errors"
	"excel-receiver/constant"
	"fmt"
	"net/http"
)

func NewF(code constant.ResCode, args string) *Error {
	return &Error{
		Code:     code,
		HttpCode: mapHttpCode(code),
		Err:      errors.New(mapMessage(code, args)),
	}
}

func mapMessage(code constant.ResCode, args string) string {
	switch code {
	case constant.InvalidRequest:
		return "invalid request"

	case constant.InvalidRowData:
		return "invalid row data"

	case constant.EmptyRowMandatory:
		return fmt.Sprintf("empty row mandatory column %s", args)

	case constant.FileSizeLimit:
		return "file size more than 128kb"

	case constant.APIUnauthorized:
		return "unauthorized"

	default:
		return "internal server error"
	}
}

func mapHttpCode(code constant.ResCode) int {
	switch code {
	case constant.InvalidRequest, constant.InvalidRowData,
		constant.EmptyRowMandatory, constant.FileSizeLimit:
		return http.StatusBadRequest

	case constant.APIUnauthorized:
		return http.StatusUnauthorized

	default:
		return http.StatusInternalServerError
	}
}

type Ierr interface {
	GetHTTPCode() int
	GetCode() string
	Error() string
	Unwrap() error
}

type Error struct {
	Code     constant.ResCode `json:"code"`
	HttpCode int
	Err      error `json:"err"`
}

func (e *Error) GetHTTPCode() int {
	return e.HttpCode
}

func (e *Error) GetCode() string {
	code := string(e.Code)

	return code
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}
