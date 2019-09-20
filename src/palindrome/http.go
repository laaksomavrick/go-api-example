package palindrome

import (
	"encoding/json"
	"net/http"
)

type apiResponse struct {
	Data   interface{} `json:"data"`
	Error  interface{}      `json:"error"`
}

// OkResponse issues a 200 http response in a uniform format across the api
func OkResponse(w http.ResponseWriter, data interface{}) {
	status := 200
	response := apiResponse{
		Data:   data,
		Error: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// ErrorResponse issues a 4xx or 5xx http response in a uniform format across the api
func ErrorResponse(w http.ResponseWriter, status int, err string) {
	response := apiResponse{
		Data: nil,
		Error:  err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}