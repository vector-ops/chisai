package cmd

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	fmt.Printf("Listening on port :8000\n")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
