package apisvr

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ResponseError struct {
	Status         int    `json:"-"`
	Name           string `json:"name"`
	Code           string `json:"code"`
	Message        string `json:"message"`
	PrivateMessage string `json:"-"` // Not sent to the client
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

func LogError(err *ResponseError) {
	if err.Name != "ServerError" {
		return // Don't log non-server errors
	}

	// .Message is not shown to public, log internally
	log.Printf("[SERVER_ERROR]: %s\n", err.PrivateMessage)
}

func ParseError(msg string) *ResponseError {
	return &ResponseError{
		Status:  http.StatusUnprocessableEntity,
		Name:    "ParseError",
		Code:    "parse_error",
		Message: msg,
	}
}

func ResourceNotFoundError(name string) *ResponseError {
	return &ResponseError{
		Status:  http.StatusNotFound,
		Name:    "ResourceNotFoundError",
		Code:    fmt.Sprintf("%s_not_found", name),
		Message: fmt.Sprintf("%s not found", name),
	}
}

func ValidationError(code string, msg string) *ResponseError {
	return &ResponseError{
		Status:  http.StatusBadRequest,
		Name:    "ValidationError",
		Code:    code,
		Message: msg,
	}
}

type FieldError struct {
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
}

type FieldErrorsMap map[string][]FieldError

func ValidationWithFieldErrors(data FieldErrorsMap) *ResponseError {
	return &ResponseError{
		Status:  http.StatusBadRequest,
		Name:    "ValidationError",
		Code:    "field_errors",
		Message: "errors in fields",
		Data:    data,
	}
}

func BindingError(err error) *ResponseError {
	var syntaxErr *json.SyntaxError
	var validationErrs validator.ValidationErrors

	switch {
	case errors.As(err, &syntaxErr):
		return ServerError(err)
	case errors.As(err, &validationErrs):
		fieldErrs := map[string][]FieldError{}
		for _, vfErr := range validationErrs {
			name := strcase.ToLowerCamel(vfErr.Field())

			fieldErr := FieldError{
				Code:    vfErr.Tag(),
				Message: vfErr.Tag(),
			}
			fieldErrs[name] = append(fieldErrs[name], fieldErr)
		}

		return ValidationWithFieldErrors(fieldErrs)
	default:
		return ServerError(err)
	}
}
