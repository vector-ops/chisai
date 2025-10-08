package routes

import (
	"net/http"

	"github.com/vector-ops/chisai/internal/api/handlers"
)

func RegisterHealthRoute(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", handlers.Health)
}
