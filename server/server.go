package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/beldpro-ci/subscriber/mailchimp"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	log "github.com/Sirupsen/logrus"
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

	if cfg.MailChimp == nil {
		err = errors.Errorf("MailChimp client can't be nil")
		return
	}

	r := mux.NewRouter()
	r.Handle("/ping",
		handlers.CombinedLoggingHandler(os.Stdout,
			http.HandlerFunc(server.PingHandler))).
		Methods("GET")

	r.Handle("/subscribe",
		handlers.CombinedLoggingHandler(os.Stdout,
			http.HandlerFunc(server.SubscribeHandler))).
		Methods("POST")

	server.mc = cfg.MailChimp
	server.router = r
	server.port = cfg.Port
	return
}

func (s *Server) Run() (err error) {
	http.Handle("/", s.router)

	err = http.ListenAndServe(fmt.Sprintf(":%d", s.port),
		handlers.CORS()(
			handlers.RecoveryHandler()(
				s.router)))
	if err != nil {
		err = errors.Wrapf(err,
			"Errored listen on port %d", s.port)
		return
	}

	return
}

func (c *Server) PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
	return
}

func isMultipart(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "multipart/form-data"
}

func isURLEncoded(r *http.Request) bool {
	return r.Header.Get("Content-Type") == "application/x-www-form-urlencoded"
}

func (c *Server) SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var contentType = r.Header.Get("Content-Type")

	switch {
	case contentType == "application/x-www-form-urlencoded":
		if err = r.ParseForm(); err != nil {
			log.
				WithError(err).
				Error("Couldn't parse form from subscribe request")
			http.Error(w,
				"Couldn't parse form",
				http.StatusBadRequest)
			return
		}
	default:
		log.Warnf("Unexpected Content-Type: %s",
			contentType)
		http.Error(w,
			"Couldn't parse form",
			http.StatusBadRequest)
		return
	}

	var email = r.FormValue("email")
	if email == "" {
		log.Warn("Tried to subscribe without email")
		http.Error(w,
			"required field 'email' not set",
			http.StatusBadRequest)
		return
	}

	if err = c.mc.Subscribe(email); err != nil {
		log.
			WithError(err).
			Error("Couldn't perform mailchimp subscription")
		http.Error(w,
			"Couldn't create MailChimp subscription",
			http.StatusInternalServerError)
		return
	}

	log.
		WithField("email", email).
		Info("User subscribed")

	fmt.Fprintf(w, "OK")
	return
}
