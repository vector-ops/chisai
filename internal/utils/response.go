package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, data any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = w.Write(payload)
	return err
}

func WriteTextResponse(w http.ResponseWriter, status int, data string) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/html")

	_, err := w.Write([]byte(data))
	return err
}

func WriteErrorResponse(w http.ResponseWriter, data error, status *int) error {
	if status == nil {
		status = toIntPointer(http.StatusInternalServerError)
	}
	return WriteJSONResponse(w, *status, map[string]any{"error": data.Error()})
}

func toIntPointer(i int) *int {
	return &i
}
