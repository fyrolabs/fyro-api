package fapi

import "net/http"

type ResponseError struct {
	Status         int    `json:"-"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Message        string `json:"message,omitempty"`
	PrivateMessage string `json:"-"`
	Data           any    `json:"data,omitempty"`
}

func (e *ResponseError) Error() string {
	return e.Message
}

func ServerError(err error) *ResponseError {
	return &ResponseError{
		Status:         http.StatusInternalServerError,
		Name:           "ServerError",
		Code:           "server_error",
		PrivateMessage: err.Error(),
	}
}
