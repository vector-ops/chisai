package middlewares

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vector-ops/chisai/internal/utils"
)

type CacheOption func(c *CacheMiddleware) error

type CacheMiddleware struct {
	rdb *redis.Client
	ttl time.Duration

	methods      []string
	excludePaths []string
	refreshKey   string
}

func (m *CacheMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		if !m.cacheablePath(r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}

		if !m.cacheableMethod(r.Method) {
			h.ServeHTTP(w, r)
			return
		}

		refresh := r.URL.Query().Has(m.refreshKey)

		if refresh {
			params := r.URL.Query()

			delete(params, m.refreshKey)

			r.URL.RawQuery = params.Encode()
		}

		sortURLParams(r.URL)

		keyI := generateKey(r.URL.String())
		if r.Method == http.MethodPost && r.Body != nil {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				h.ServeHTTP(w, r)
				return
			}

			rd := io.NopCloser(bytes.NewBuffer(body))
			keyI = generateKeyWithBody(r.URL.String(), body)

			r.Body = rd
		}

		key := strconv.FormatUint(keyI, 10)

		if refresh {
			m.rdb.Del(r.Context(), key)
		} else {
			err := m.readCache(r.Context(), key, w)
			if err == nil {
				slog.Info("Cache hit", "time", time.Since(start).String())
				return
			}
		}

		recorder := utils.NewResponseRecorder(w)

		h.ServeHTTP(recorder, r)

		if err := m.cacheResult(r.Context(), key, recorder); err != nil {
			slog.Error("Failed to cache response", "error", err)
		}
		slog.Info("Cache miss", "time", time.Since(start).String())
	})

}

func (m *CacheMiddleware) readCache(ctx context.Context, key string, w http.ResponseWriter) error {

	value, err := m.rdb.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	entry := &utils.CacheEntry{}
	if err := entry.Decode(bytes.NewBufferString(value).Bytes()); err != nil {
		return err
	}

	return entry.Replay(w)
}

func (m *CacheMiddleware) cacheResult(ctx context.Context, key string, r *utils.ResponseRecorder) error {
	e := r.Result()

	if e.StatusCode >= 400 {
		return nil
	}

	b, err := e.Encode()
	if err != nil {
		return fmt.Errorf("unable to read recorded response: %s", err)
	}

	expiration := m.ttl + time.Duration(rand.Intn(10))*time.Second

	_, err = m.rdb.Set(
		ctx,
		key,
		string(b),
		expiration,
	).Result()

	return err
}

func generateKey(URL string) uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(URL))

	return hash.Sum64()
}

func generateKeyWithBody(URL string, body []byte) uint64 {
	hash := fnv.New64a()
	body = append([]byte(URL), body...)
	hash.Write(body)

	return hash.Sum64()
}

func sortURLParams(URL *url.URL) {
	params := URL.Query()

	for _, param := range params {
		slices.Sort(param)
	}

	URL.RawQuery = params.Encode()
}

func (m *CacheMiddleware) cacheablePath(path string) bool {
	return !slices.Contains(m.excludePaths, path)
}

func (m *CacheMiddleware) cacheableMethod(method string) bool {
	return slices.Contains(m.methods, method)
}

func NewCacheMiddleware(opts ...CacheOption) (*CacheMiddleware, error) {

	c := &CacheMiddleware{}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	if c.rdb == nil {
		return nil, errors.New("cache middleware requires redis client")
	}

	if int64(c.ttl) < 1 {
		return nil, errors.New("cache ttl is not set")
	}

	if c.methods == nil {
		c.methods = []string{"GET"}
	}

	return c, nil

}

func WithRedisClient(rdb *redis.Client) CacheOption {
	return func(c *CacheMiddleware) error {
		c.rdb = rdb

		return nil
	}
}

func WithTTL(ttl time.Duration) CacheOption {
	return func(c *CacheMiddleware) error {

		if int64(ttl) < 1 {
			return fmt.Errorf("cache middleware ttl %v is invalid", ttl)
		}

		c.ttl = ttl

		return nil
	}
}

func WithMethods(methods []string) CacheOption {
	return func(c *CacheMiddleware) error {
		for _, method := range methods {
			if method != http.MethodGet && method != http.MethodPost {
				return fmt.Errorf("invalid method %s", method)
			}
		}
		c.methods = methods
		return nil
	}
}

func WithExcludePaths(paths []string) CacheOption {
	return func(c *CacheMiddleware) error {
		c.excludePaths = paths
		return nil
	}
}

func WithRefreshKey(refreshKey string) CacheOption {
	return func(c *CacheMiddleware) error {
		c.refreshKey = refreshKey
		return nil
	}
}
