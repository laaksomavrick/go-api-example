package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

const (
	GET    = http.MethodGet
	POST   = http.MethodPost
	PATCH  = http.MethodPatch
	DELETE = http.MethodDelete
)

// Server defines the shape of the server, ie, the dependencies this service requires
// in order to function. This is the top-level object which all handlers have injected into them, and thus
// provides shared resources (e.g database) to all downstream components (e.g repositories). This application
// architecture provides some nice properties (easily testable components, easy mocking, no DI framework magic)
// at the cost of a little boilerplate (see main.go) and set up (see wire()).
type Server struct {
	router *mux.Router
	db     *sqlx.DB
	config *Config
}

// ServerHandlerFunc defines the shape of a http handler for use with the Server.
type ServerHandlerFunc = func(server *Server) http.HandlerFunc

// NewServer constructs a new server struct.
func NewServer(router *mux.Router, db *sqlx.DB, config *Config) *Server {
	return &Server{
		router: router,
		db:     db,
		config: config,
	}
}

// Init initializes the server, mapping routes to their handlers and serving the http server.
func (s *Server) Init() {
	s.wire()
	s.serve()
}

// Wire sets up the http handlers and maps routes to their handlers
func (s *Server) wire() {
	// In a real API, we'd probably have some middlewares we'd like to apply for common functionality
	// e.g authentication and authorization; logging requests
	// See https://github.com/gorilla/mux#middleware

	// DI the server object to our handlers
	// this makes mocking dependencies easy to do in tests and has the pleasant property
	// of making available useful, shared functionality (ie: database, config)
	getHealthzHandler := s.inject(GetHealthzHandler)

	getMessagesHandler := s.inject(GetMessagesHandler)
	postMessageHandler := s.inject(PostMessageHandler)
	getMessageHandler := s.inject(GetMessageHandler)
	patchMessageHandler := s.inject(PatchMessageHandler)
	deleteMessageHandler := s.inject(DeleteMessageHandler)

	// Though not necessary, a /healthz endpoint typically exists for checking the service health
	// Additionally, if we were using something like prometheus, we'd have a /metrics as well
	s.router.HandleFunc("/healthz", getHealthzHandler).Methods(GET)

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
