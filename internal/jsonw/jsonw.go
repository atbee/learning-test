package jsonw

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Msg    string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.Status, e.Msg)
}

func NewErrorResponse(status string, msg string) error {
	return &ErrorResponse{
		Status: status,
		Msg:    msg,
	}
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(NewErrorResponse("Internal Server Error", err.Error()))
}

func Unauthorized(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(NewErrorResponse("Unauthorized", err.Error()))
}

func BadRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(NewErrorResponse("Bad Request", err.Error()))
}

func Forbidden(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(NewErrorResponse("Forbidden", err.Error()))
}
