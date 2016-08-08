package redis

import (
	"fmt"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	caddy.RegisterPlugin("redis", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	// Parse the arguments
	server, password, err := parse(c)
	if err != nil {
		return err
	}
	// Set the server and password for the redis pool
	pool = newPool(server, password)

	cfg := httpserver.GetConfig(c)
	mid := func(next httpserver.Handler) httpserver.Handler {
		return Redis{Next: next}
	}
	cfg.AddMiddleware(mid)

	return nil
}

func parse(c *caddy.Controller) (string, string, error) {
	server := ":6379"
	password := ""

	for c.Next() {
		for c.NextBlock() {
			switch c.Val() {
			case "server":
				if !c.NextArg() {
					return server, password, c.ArgErr()
				}
				server = c.Val()
				fmt.Println(LogTag, "Found server config:", server)
			case "password":
				if !c.NextArg() {
					return server, password, c.ArgErr()
				}
				password = c.Val()
				fmt.Println(LogTag, "Found password config:", password)
			}
		}
	}

	return server, password, nil
}
