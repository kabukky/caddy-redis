# caddy-redis

This package is a plugin for the [Caddy](https://caddyserver.com) webserver. If you POST valid JSON to `/redis/<key>`, it will put it into into a Redis database.

Similarly, you do a GET on `/redis/<key>`, it will return the previously stored data. If you specify a `Accept-Encoding` header of `application/xml` for the GET request, it will return the data as XML, otherwise JSON.

Sample usage: 

`curl -X POST -H "Content-Type: application/json" -d '{"key":"val"}' http://localhost:2015/redis/test`

This will SET `caddy:test => "{"key":"val"}"` in Redis.

`curl -H "Accept-Encoding: application/json" http://localhost:2020/redis/test`

This will GET the previously stored data for the key `test` as JSON.

`curl -H "Accept-Encoding: application/xml" http://localhost:2020/redis/test`

This will GET the previously stored data for the key `test` as XML.

**WARNING:** This plugin is for testing purposes only. **Do not use in production!** Seriously, this plugin will write anything given by anyone to the specified Redis database. You'll want to implement some kind of authentication at the very least.

## Syntax

```
redis  {
	server localhost:6379 	# NOTE: If not specified, ":6379" is used as the Redis server (default).
	password foobar 		# NOTE: If not specified, Redis is used without authentication (default).
}
```
