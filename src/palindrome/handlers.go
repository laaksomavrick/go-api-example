package palindrome

import (
	"net/http"
)

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	OkResponse(w, map[string]interface{}{"hello": "world"})
}