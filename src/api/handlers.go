package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	statusInternalServerErrorMessage = "Internal server error"
	statusBadRequestMessage          = "Bad request"
	statusNotFoundMessage            = "Requested resource not found"
)

// GetHealthzHandler handles incoming requests against the healthz endpoint.
func GetHealthzHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// This might do some pings to other processes which this service depends on e.g postgres
		OkResponse(w, "ok")
	}
}

// GetMessagesHandler handles incoming requests against the get messages endpoint.
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

// PostMessageHandler handles incoming requests against the create message endpoint.
func PostMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upsertMessageDto := UpsertMessageDto{}
		repo := NewRepository(s.db)

		// Map body to dto
		err := json.NewDecoder(r.Body).Decode(&upsertMessageDto)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		// Validate dto is sound
		err = upsertMessageDto.Validate()

		if err != nil {
			ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Create message
		message, err := repo.CreateMessage(upsertMessageDto.Content)

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, message)
	}
}

// GetMessageHandler handles incoming requests against the get message endpoint.
func GetMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)

		// Parse {id} from url
		id, err := getIdFromUrl(r)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		// Get the message
		message, err := repo.GetMessage(id)

		if err != nil {
			ErrorResponse(w, http.StatusNotFound, statusNotFoundMessage)
			return
		}

		OkResponse(w, message)
	}
}

// PatchMessageHandler handles incoming requests against the update message endpoint.
func PatchMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		upsertMessageDto := UpsertMessageDto{}
		repo := NewRepository(s.db)

		// Parse {id} from url
		id, err := getIdFromUrl(r)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		// Map body to dto
		err = json.NewDecoder(r.Body).Decode(&upsertMessageDto)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		// Validate dto is sound
		err = upsertMessageDto.Validate()

		if err != nil {
			ErrorResponse(w, http.StatusUnprocessableEntity, err.Error())
			return
		}

		// Check that message exists
		_, err = repo.GetMessage(id)

		if err != nil {
			ErrorResponse(w, http.StatusNotFound, statusNotFoundMessage)
			return
		}

		// Update the message
		message, err := repo.UpdateMessage(id, upsertMessageDto.Content)

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, message)
	}
}

// DeleteMessageHandler handles incoming requests against the delete message endpoint.
func DeleteMessageHandler(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repo := NewRepository(s.db)

		// Parse {id} from url
		id, err := getIdFromUrl(r)

		if err != nil {
			ErrorResponse(w, http.StatusBadRequest, statusBadRequestMessage)
			return
		}

		// Check that message exists
		_, err = repo.GetMessage(id)

		if err != nil {
			ErrorResponse(w, http.StatusNotFound, statusNotFoundMessage)
			return
		}

		// Delete the message
		err = repo.DeleteMessage(id)

		if err != nil {
			ErrorResponse(w, http.StatusInternalServerError, statusInternalServerErrorMessage)
			return
		}

		OkResponse(w, nil)
	}
}

// getIdFromUrl is a utility function to parse the {id} parameter out of a given request's url
func getIdFromUrl(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	maybeId := vars["id"]
	id, err := strconv.Atoi(maybeId)
	return id, err
}
