package routes

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/vector-ops/chisai/internal/api/handlers"
)

func RegisterURLRoutes(mux *http.ServeMux, db *dynamodb.Client) {
	h := handlers.NewUrlHandler(db)

	mux.HandleFunc("GET /", h.GetUrl)
	mux.HandleFunc("POST /shorten", h.Shorten)
}
