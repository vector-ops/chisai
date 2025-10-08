package handlers

import (
	"net/http"

	"github.com/vector-ops/chisai/internal/utils"
)

func Health(w http.ResponseWriter, r *http.Request) {
	// utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"message": "Hello, World!"})
	utils.WriteTextResponse(w, http.StatusOK, "Hello, !")
}
