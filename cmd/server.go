package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/vector-ops/chisai/internal/api/routes"
)

type Server struct {
	mux *http.ServeMux
	db  *dynamodb.Client
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	err := InitDB()
	if err != nil {
		log.Fatal(err)
	}
	db := GetDB()

	routes.RegisterHealthRoute(mux)
	routes.RegisterURLRoutes(mux, db)

	fmt.Printf("Listening on port :8000\n")
	log.Fatal(http.ListenAndServe(":8000", mux))
}
