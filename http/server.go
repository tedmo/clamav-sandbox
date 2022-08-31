package http

import (
	"context"
	"errors"
	"github.com/tedmo/cav/logger"
	"io"
	"net/http"
)

type Scanner interface {
	Scan(ctx context.Context, contents io.Reader) ([]byte, error)
}

type Server struct {
	Scanner Scanner
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.routes())
}

func (s *Server) routes() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := s.handleFile(w, r); err != nil {
			logger.New(r.Context()).Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return router
}

func (s *Server) handleFile(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		return errors.New("unsupported method")
	}
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return err
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return err
	}

	filename := header.Filename
	if filename == "" {
		return errors.New("filename must not be empty")
	}

	status := http.StatusOK
	output, err := s.Scanner.Scan(r.Context(), file)
	if err != nil {
		status = http.StatusInternalServerError
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(status)
	w.Write(output)

	return nil
}
