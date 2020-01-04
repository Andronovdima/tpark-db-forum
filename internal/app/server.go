package apiserver

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"net/http"
)

type Server struct {
	Mux          *mux.Router
	SessionStore sessions.Store
	Config       *Config
	Logger       *zap.SugaredLogger
}

func NewServer(config *Config, logger *zap.SugaredLogger) (*Server, error) {
	s := &Server{
		Mux:          mux.NewRouter(),
		SessionStore: sessions.NewCookieStore([]byte(config.SessionKey)),
		Logger:       logger,
		Config:       config,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Mux.ServeHTTP(w, r)
}

func (s *Server) ConfigureServer(db *sql.DB) {
	//
}
