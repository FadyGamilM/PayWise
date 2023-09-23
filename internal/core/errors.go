package core

import (
	"fmt"
	"net/http"
)

type AppErr string

const (
	NOT_FOUND_ERR       AppErr = "resource not found"
	BAD_REQUEST_ERR     AppErr = "bad request"
	UNAUTHORIZED_ERR    AppErr = "not authorized"
	INTERNAL_SERVER_ERR AppErr = "internal server error"
)

type ERROR struct {
	Type AppErr `json:"error_type"`
	Msg  string `json:"error_msg"`
}

func (e *ERROR) Error() string {
	return e.Msg
}

func (e *ERROR) StatusCode() int {
	switch e.Type {
	case UNAUTHORIZED_ERR:
		return http.StatusUnauthorized
	case INTERNAL_SERVER_ERR:
		return http.StatusInternalServerError
	case BAD_REQUEST_ERR:
		return http.StatusBadRequest
	case NOT_FOUND_ERR:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func NewNotFoundError(resource string, idv string) ERROR {
	return ERROR{
		Type: NOT_FOUND_ERR,
		Msg:  fmt.Sprintf("no %v found with key : %v", resource, idv),
	}
}

//	func NewBadRequestError(fields ...string) ERROR {
//		var errorMsg strings.Builder
//		for _, errField := range fields {
//			errorMsg.WriteString("field " + errField + " is wrong\n")
//		}
//		return ERROR{
//			Type: BAD_REQUEST_ERR,
//			Msg:  errorMsg.String(),
//		}
//	}
func NewBadRequestError() ERROR {
	return ERROR{
		Type: BAD_REQUEST_ERR,
		Msg:  "valdiation errors",
	}
}

func NewUnAuthorizedError() ERROR {
	return ERROR{
		Type: UNAUTHORIZED_ERR,
		Msg:  "not authorized to access this resource or perform this action",
	}
}

func NewInternalServerError() ERROR {
	return ERROR{
		Type: INTERNAL_SERVER_ERR,
		Msg:  "something went wrong, internal server error",
	}
}
