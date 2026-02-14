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

func BadRequest(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusBadRequest, errMsg)
}

func NotFound(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusNotFound, errMsg)
}

func Internal(w http.ResponseWriter, errMsg string) {
	Err(w, http.StatusInternalServerError, errMsg)
}
