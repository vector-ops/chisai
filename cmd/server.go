package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vector-ops/chisai/internal/api/routes"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	mux *http.ServeMux
	db  *mongo.Database
}

func NewServer(db *mongo.Database) *Server {
	return &Server{db: db}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	routes.RegisterHealthRoute(mux)
	routes.RegisterURLRoutes(mux, s.db)

	fmt.Printf("Listening on port :8000\n")
	log.Fatal(http.ListenAndServe(":8000", mux))
}

func (s *Server) Close() error {
	return nil
}
