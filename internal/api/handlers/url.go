package handlers

import (
	"net/http"
	"strings"

	"github.com/vector-ops/chisai/internal/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type URLStore struct {
	ShortURL    string `json:"shortUrl" bson:"shortUrl"`
	RedirectURL string `json:"redirectUrl" bson:"redirectUrl"`
}

type ShortenURLRequestPayload struct {
	URL string `json:"url"`
}

type URLHandler struct {
	db *mongo.Database
}

func NewURLHandler(db *mongo.Database) *URLHandler {
	return &URLHandler{
		db: db,
	}
}

func (h *URLHandler) Shorten(w http.ResponseWriter, r *http.Request) {
	var payload ShortenURLRequestPayload

	err := utils.ReadJSONBody(r, &payload)
	if err != nil {
		utils.WriteErrorResponse(w, err, nil)
		return
	}

	if payload.URL == "" {
		utils.WriteErrorResponse(w, utils.ErrMissingURL, nil)
		return
	}

	shortURL := utils.GenerateShortURL(payload.URL)

	item := URLStore{ShortURL: shortURL, RedirectURL: payload.URL}

	doc := bson.M{"$set": item}

	opts := options.UpdateOne().SetUpsert(true)

	_, err = h.db.Collection("url_store").UpdateOne(r.Context(), bson.M{"redirectUrl": payload.URL}, doc, opts)
	if err != nil {
		utils.WriteErrorResponse(w, err, nil)
		return
	}

	utils.WriteJSONResponse(w, 201, map[string]any{"shortUrl": shortURL})
}

func (h *URLHandler) GetURL(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	shortURL := strings.TrimLeft(p, "/")

	var doc URLStore
	err := h.db.Collection("url_store").FindOne(r.Context(), bson.M{"shortUrl": shortURL}).Decode(&doc)
	if err != nil {
		utils.WriteErrorResponse(w, err, nil)
		return
	}

	http.Redirect(w, r, doc.RedirectURL, http.StatusMovedPermanently)
}
