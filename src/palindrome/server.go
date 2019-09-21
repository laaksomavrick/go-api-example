package palindrome

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

const (
	GET = http.MethodGet
	POST = http.MethodPost
	PATCH = http.MethodPatch
	DELETE = http.MethodDelete
)

type Server struct {
	router *mux.Router
	db *sqlx.DB
	config *Config
}

type ServerHandlerFunc = func(server *Server) http.HandlerFunc

func NewServer(router *mux.Router, db *sqlx.DB, config *Config) *Server {
	return &Server{
		router: router,
		db: db,
		config: config,
	}
}

func (s *Server) Init() {
	s.wire()
	s.serve()
}

func (s *Server) wire() {
	// In a real API, we'd probably have some middlewares we'd like to apply for common functionality
	// e.g authentication and authorization; logging requests
	// See https://github.com/gorilla/mux#middleware

	// DI the server object to our handlers
	// this makes mocking dependencies easy to do in tests and has the pleasant property
	// of making available useful, shared functionality (ie: database, config)
	getMessagesHandler := s.inject(GetMessagesHandler)
	postMessageHandler := s.inject(PostMessageHandler)
	getMessageHandler := s.inject(GetMessageHandler)
	patchMessageHandler := s.inject(PatchMessageHandler)
	deleteMessageHandler := s.inject(DeleteMessageHandler)

	s.router.HandleFunc("/messages", getMessagesHandler).Methods(GET)
	s.router.HandleFunc("/message", postMessageHandler).Methods(POST)
	s.router.HandleFunc("/message/{id}", getMessageHandler).Methods(GET)
	s.router.HandleFunc("/message/{id}", patchMessageHandler).Methods(PATCH)
	s.router.HandleFunc("/message/{id}", deleteMessageHandler).Methods(DELETE)

	// We want clients to be able to talk to our api (e.g a single page app, a CLI tool, whatever)
	s.router.Use(mux.CORSMethodMiddleware(s.router))
}

func (s *Server) serve() {
	port := fmt.Sprintf(":%s", s.config.port)
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(port, s.router))
}

func (s *Server) inject(handlerFunc ServerHandlerFunc) http.HandlerFunc {
	return handlerFunc(s)
}
