package palindrome

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	statusInternalServerErrorMessage = "Internal server error"
	statusBadRequestMessage = "Bad request"
	statusUnprocessableEntityMessage = "Unprocessable entity"
)

func GetMessagesHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)

		// In a production API, this would probably need to implement pagination using SIZE and OFFSET
		// so we don't return too many rows
		messages, err := repo.GetMessages()

		if err != nil {
			log.Print(err)
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, messages)
	}
}

func PostMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createMessageDto := CreateMessageDto{}
		repo := NewRepository(s.db)

		err := json.NewDecoder(r.Body).Decode(&createMessageDto)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		err = createMessageDto.Validate()

		if err != nil {
			ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		message, err := repo.CreateMessage(createMessageDto.Content)

		if err != nil {
			log.Print(err)
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, message)
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
