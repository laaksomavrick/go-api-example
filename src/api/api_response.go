package api

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Resource   interface{} `json:"resource"`
	Error  interface{}      `json:"error"`
}

// OkResponse issues a 200 http response in a uniform format across the api
func OkResponse(w http.ResponseWriter, data interface{}) {
	status := 200
	response := apiResponse{
		Resource:   data,
		Error: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

// ErrorResponse issues a 4xx or 5xx http response in a uniform format across the api
func ErrorResponse(w http.ResponseWriter, status int, err string) {
	response := apiResponse{
		Resource: nil,
		Error:  err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}