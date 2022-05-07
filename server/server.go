package server

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
)

type Server struct {
	wg  sync.WaitGroup
	mux *http.ServeMux
}

func New() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) WithMetrics() *Server {
	s.mux.Handle("/metrics", promhttp.Handler())
	return s
}

func (s *Server) WithHandler(url string, handler func(http.ResponseWriter, *http.Request)) *Server {
	s.mux.HandleFunc(url, handler)
	return s
}

func (s *Server) Start() {
	srv := http.Server{
		Addr:    ":8000",
		Handler: s.mux,
	}

	go func() {
		s.wg.Add(1)
		defer s.wg.Done()
		err := srv.ListenAndServe()
		if err != nil {
			return
		}
	}()
}

func (s *Server) Stop() {
	s.wg.Wait()
}
