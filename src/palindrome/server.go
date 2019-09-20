package palindrome

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type Server struct {
	router *mux.Router
	db *sqlx.DB
	config *Config
}

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
	s.router.HandleFunc("/", HelloWorldHandler)
}

func (s *Server) serve() {
	port := fmt.Sprintf(":%s", s.config.port)
	log.Printf("listening on %s", port)
	log.Fatal(http.ListenAndServe(port, s.router))
}
