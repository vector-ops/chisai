package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadJSONBody(r *http.Request, dest interface{}) (err error) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	err = json.Unmarshal(body, dest)

	return
}
