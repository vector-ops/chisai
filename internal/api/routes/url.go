package routes

import (
	"net/http"

	"github.com/vector-ops/chisai/internal/api/handlers"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterURLRoutes(mux *http.ServeMux, db *mongo.Database) {
	h := handlers.NewURLHandler(db)

	mux.HandleFunc("GET /", h.GetURL)
	mux.HandleFunc("POST /shorten", h.Shorten)
}
