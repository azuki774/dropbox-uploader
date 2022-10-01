package server

import (
	"fmt"
	"net/http"
	"time"
)

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "It is the root page.\n")
}

func (s *Server) uploadHandler(w http.ResponseWriter, r *http.Request) {
	const uploadfilekey = "file"
	file, _, err := r.FormFile(uploadfilekey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}
	defer file.Close()

	q := r.URL.Query()
	path := q.Get("path")
	err = s.Usecase.UploadFile(path, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error: %s\n", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}
