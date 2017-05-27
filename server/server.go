package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	port   int
	router *mux.Router
}

type Config struct {
	Port int
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func New(cfg Config) (server Server, err error) {
	if cfg.Port == 0 {
		err = errors.Errorf("port can't be 0")
		return
	}

	r := mux.NewRouter()
	r.Handle("/ping", handlers.CombinedLoggingHandler(os.Stdout,
		http.HandlerFunc(pingHandler)))
	r.Handle("/subscribe", handlers.CombinedLoggingHandler(os.Stdout,
		http.HandlerFunc(subscribeHandler)))

	server.router = r
	server.port = cfg.Port
	return
}

func (s *Server) Run() (err error) {
	http.Handle("/", s.router)

	err = http.ListenAndServe(fmt.Sprintf(":%d", s.port),
		handlers.RecoveryHandler()(s.router))
	if err != nil {
		err = errors.Wrapf(err,
			"Errored listen on port %d", s.port)
		return
	}

	return
}
