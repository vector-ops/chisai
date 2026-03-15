package cmd

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vector-ops/chisai/internal/api/handlers"
	"github.com/vector-ops/chisai/internal/api/middlewares"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Server struct {
	mux *http.ServeMux
	db  *mongo.Database
	rdb *redis.Client
}

func NewServer(db *mongo.Database, rdb *redis.Client) *Server {
	return &Server{db: db, rdb: rdb}
}

func (s *Server) Start() {
	mux := http.NewServeMux()

	cacheMiddleware, err := middlewares.NewCacheMiddleware(
		middlewares.WithRedisClient(s.rdb),
		middlewares.WithTTL(time.Minute*5),
		middlewares.WithExcludePaths([]string{"/shorten", "/health"}),
		middlewares.WithRefreshKey("refresh"),
	)

	if err != nil {
		slog.Error("Failed to initialize cache middleware", "error", err)
		return
	}

	mux.HandleFunc("GET /health", handlers.Health)

	h := handlers.NewURLHandler(db)
	mux.HandleFunc("GET /", h.GetURL)
	mux.HandleFunc("POST /shorten", h.Shorten)

	slog.Info("Listening on port :8000")
	log.Fatal(http.ListenAndServe(":8000", cacheMiddleware.Handler(mux)))
}

func (s *Server) Close() error {
	return nil
}
