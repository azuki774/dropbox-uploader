package server

import (
	"azuki774/dropbox-uploader/internal/usecases"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	Host    string
	Port    string
	Logger  *zap.Logger
	Usecase *usecases.Usecases
}

func (s *Server) Start() error {
	router := mux.NewRouter()
	s.addRecordFunc(router)
	return http.ListenAndServe(net.JoinHostPort(s.Host, s.Port), router)
}

func (s *Server) addRecordFunc(r *mux.Router) {
	r.HandleFunc("/", s.rootHandler)
	r.HandleFunc("/upload", s.uploadHandler).Queries("path", "{path}").Methods("POST")
	r.Use(s.middlewareLogging)
}

func (s *Server) middlewareLogging(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			s.Logger.Info("access", zap.String("url", r.URL.Path), zap.String("X-Forwarded-For", r.Header.Get("X-Forwarded-For")))
		}
		h.ServeHTTP(w, r)
	})
}
