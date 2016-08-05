package redis

import (
	"net/http"

	"github.com/mholt/caddy/caddyhttp/httpserver"
)

var LOG_TAG = "caddy-redis"

type Redis struct {
	Next httpserver.Handler
}

func (redis Redis) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	if r.URL.Path == "/redis" || r.URL.Path == "/redis/" {
		query := r.URL.Query()
		for key, value := range query {
			set(key, value[0])
		}
	}
	return redis.Next.ServeHTTP(w, r)
}
