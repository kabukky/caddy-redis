# caddy-redis

This package is a plugin for the [Caddy](https://caddyserver.com) webserver. It will write any field-value pair specified in the query parameters of a HTTP request to `/redis` to a Redis database.

Sample usage: 

`http://localhost:2015/redis?test=1234&foo=bar`

This will SET `caddy:test => "1234"` and `caddy:foo => "bar"` in Redis.

**WARNING:** This plugin is for testing purposes only. Do not use in production! Seriously, this plugin will write anything given by anyone to the specified Redis database. You'll want to implement some kind of authentication at the very least.

## Syntax

```
redis  {
	server localhost:6379 # NOTE: If not specified, ":6379" is used as the Redis server (default).
	password foobar # NOTE: If not specified, Redis is used without authentication (default).
}
```
