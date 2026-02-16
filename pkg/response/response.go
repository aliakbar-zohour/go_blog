// pkg/response: Consistent JSON response format for the API.
package response

import (
	"encoding/json"
	"net/http"
)

type Body struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"` // Machine-readable error code for clients
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func OK(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Body{Success: true, Data: data})
}

func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Body{Success: true, Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Err(w http.ResponseWriter, status int, errMsg string) {
	JSON(w, status, Body{Success: false, Error: errMsg})
}

// ErrWithCode sends an error response with a machine-readable code (e.g. invalid_email, not_found).
func ErrWithCode(w http.ResponseWriter, status int, code, errMsg string) {
	JSON(w, status, Body{Success: false, Error: errMsg, Code: code})
}

func BadRequest(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusBadRequest, errMsg)
}

func NotFound(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusNotFound, errMsg)
}

func Internal(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusInternalServerError, errMsg)
}

func Unauthorized(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusUnauthorized, errMsg)
}

func Forbidden(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusForbidden, errMsg)
}

// BadRequestWithCode sends 400 with an error code.
func BadRequestWithCode(w http.ResponseWriter, code, errMsg string) {
	ErrWithCode(w, http.StatusBadRequest, code, errMsg)
}

// UnauthorizedWithCode sends 401 with an error code.
func UnauthorizedWithCode(w http.ResponseWriter, code, errMsg string) {
	ErrWithCode(w, http.StatusUnauthorized, code, errMsg)
}

// NotFoundWithCode sends 404 with an error code.
func NotFoundWithCode(w http.ResponseWriter, code, errMsg string) {
	ErrWithCode(w, http.StatusNotFound, code, errMsg)
}
