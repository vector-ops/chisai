package handlers

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/vector-ops/chisai/internal/utils"
)

type ShortenURLRequestPayload struct {
	URL string `json:"url"`
}

type UrlHandler struct {
	db *dynamodb.Client
}

func NewUrlHandler(db *dynamodb.Client) *UrlHandler {
	return &UrlHandler{
		db: db,
	}
}

func (h *UrlHandler) Shorten(w http.ResponseWriter, r *http.Request) {
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

	shortUrl := utils.GenerateShortURL(payload.URL)

	item, err := attributevalue.MarshalMap(map[string]string{"shortUrl": shortUrl})
	if err != nil {
		utils.WriteErrorResponse(w, err, nil)
		return
	}

	_, err = h.db.PutItem(r.Context(), &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("url_store"),
	})
	if err != nil {
		utils.WriteErrorResponse(w, err, nil)
		return
	}

	utils.WriteJSONResponse(w, 201, map[string]any{"shortUrl": shortUrl})
}

func (h *UrlHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	// h.db.GetItem(r.Context(), &dynamodb.GetItemInput{
	// 	Key: map[string]types.AttributeValue{
	// 		"shortUrl": attributevalue.Marshal(payload.URL),
	// 	},
	// 	TableName: aws.String("url_store"),
	// })
	utils.WriteJSONResponse(w, 301, map[string]any{"message": p})
}
