package palindrome

import (
	"net/http"
)

func GetMessagesHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)
		messages, err := repo.GetMessages()

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, "Oops! Something went wrong.")
			return
		}

		OkResponse(w, messages)
	}
}

func PostMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate request
		// Call isPalindrome
		// Insert
		// Return newly created record
		OkResponse(w, map[string]interface{}{"hello": "world"})
	}
}

func GetMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check that id exists
		// Get the message
		// Return it
		OkResponse(w, map[string]interface{}{"hello": "world"})
	}
}

func PatchMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check that id exists
		// Check that request is valid
		// Update the message
		// Return the updated message
		OkResponse(w, map[string]interface{}{"hello": "world"})
	}
}

func DeleteMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check that id exists
		// Delete the resource
		// Return nothing
		OkResponse(w, map[string]interface{}{"hello": "world"})
	}
}
