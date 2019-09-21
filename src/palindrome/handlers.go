package palindrome

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	statusInternalServerErrorMessage = "Internal server error"
	statusBadRequestMessage = "Bad request"
	statusUnprocessableEntityMessage = "Unprocessable entity"
	statusNotFoundMessage = "Requested resource not found"
)

func GetMessagesHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)

		// In a production API, this would probably need to implement pagination using SIZE and OFFSET
		// so we don't return too many rows
		messages, err := repo.GetMessages()

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, messages)
	}
}

func PostMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upsertMessageDto := UpsertMessageDto{}
		repo := NewRepository(s.db)

		err := json.NewDecoder(r.Body).Decode(&upsertMessageDto)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		err = upsertMessageDto.Validate()

		if err != nil {
			ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		message, err := repo.CreateMessage(upsertMessageDto.Content)

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, message)
	}
}

func GetMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)
		id, err := getIdFromUrl(r)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		message, err := repo.GetMessage(id)

		if err != nil {
			ErrorResponse(w, http.StatusNotFound, statusNotFoundMessage)
			return
		}

		OkResponse(w, message)
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

func getIdFromUrl(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	maybeId := vars["id"]
	id, err := strconv.Atoi(maybeId)
	return id, err
}