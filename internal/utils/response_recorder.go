package utils

import (
	"bytes"
	"encoding/gob"
	"maps"
	"net/http"
)

type ResponseRecorder struct {
	http.ResponseWriter

	status       int
	body         bytes.Buffer
	headers      http.Header
	headerCopied bool
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		headers:        make(http.Header),
	}
}

func (w *ResponseRecorder) Write(b []byte) (int, error) {
	w.copyHeaders()
	i, err := w.ResponseWriter.Write(b)
	if err != nil {
		return i, err
	}

	return w.body.Write(b[:i])
}

func (r *ResponseRecorder) copyHeaders() {
	if r.headerCopied {
		return
	}

	r.headerCopied = true
	maps.Copy(r.headers, r.Header())
}

func (w *ResponseRecorder) WriteHeader(statusCode int) {
	w.copyHeaders()

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Result() *CacheEntry {
	r.copyHeaders()

	return &CacheEntry{
		Header:     r.headers,
		StatusCode: r.status,
		Body:       r.body.Bytes(),
	}
}

type CacheEntry struct {
	Header     http.Header
	StatusCode int
	Body       []byte
}

func (c *CacheEntry) Encode() ([]byte, error) {
	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(c)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *CacheEntry) Decode(b []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(b))
	return dec.Decode(c)
}

func (c *CacheEntry) Replay(w http.ResponseWriter) error {
	maps.Copy(w.Header(), c.Header)

	if c.StatusCode != 0 {
		w.WriteHeader(c.StatusCode)
	}

	if len(c.Body) == 0 {
		return nil
	}

	_, err := w.Write(c.Body)
	return err
}
