package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type Server struct {
	port   int
	mc     *mailchimp.Client
	router *mux.Router
}

type Config struct {
	MailChimp *mailchimp.Client
	Port      int
}

func New(cfg Config) (server Server, err error) {
	if cfg.Port == 0 {
		err = errors.Errorf("port can't be 0")
		return
	}

	r := mux.NewRouter()
	r.Handle("/ping", handlers.CombinedLoggingHandler(os.Stdout,
		http.HandlerFunc(server.pingHandler)))
	r.Handle("/subscribe", handlers.CombinedLoggingHandler(os.Stdout,
		http.HandlerFunc(server.subscribeHandler)))

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

func (c *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
	return
}

func (c *Server) subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w,
			"Couldn't parse form",
			http.StatusBadRequest)
		return
	}

	var email = r.FormValue("email")
	if email == "" {
		http.Error(w,
			"required field 'email' not set",
			http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "OK")
	return
}
