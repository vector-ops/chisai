package utils

import "encoding/base64"

func GenerateShortURL(url string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(url))
}
