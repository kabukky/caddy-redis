package redis

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/clbanning/mxj"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

var (
	LogTag       = "caddy-redis"
	apiPath      = "/redis/"
	apiPathRegex = regexp.MustCompile(apiPath + "(.+[^/])/?")
)

type Redis struct {
	Next httpserver.Handler
}

func (redis Redis) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {
	// Get key from url
	match := apiPathRegex.FindStringSubmatch(r.URL.Path)

	// See if we got a match
	if len(match) > 1 {
		key := match[1]
		// Check HTTP method
		if r.Method == http.MethodPost {
			// POST
			// We need to write to the database
			// Get json
			// Here we are using simplejson to convert to a json object. This is just to check if we were given valid json with the request
			json, err := simplejson.NewFromReader(r.Body)
			if err != nil {
				errorMessage := "Error while decoding json from POST: " + err.Error()
				fmt.Println(LogTag, errorMessage)
				serveJSONError(w, errorMessage)
				return 0, nil
			}
			fmt.Println(LogTag, "Got json from POST:", json)
			jsonEncoded, err := json.Encode()
			if err != nil {
				errorMessage := "Error while encoding json from POST: " + err.Error()
				fmt.Println(LogTag, errorMessage)
				serveJSONError(w, errorMessage)
				return 0, nil
			}
			// Write the json string to the databse under the key
			set(key, string(jsonEncoded))
		} else if r.Method == http.MethodGet {
			// GET
			// We need to get the key from the database
			value, err := get(key)
			if err != nil {
				errorMessage := "Error while getting value from Redis: " + err.Error()
				fmt.Println(LogTag, errorMessage)
				serveJSONError(w, errorMessage)
				return 0, nil
			}
			// Here we are using simplejson to convert to a json object. This is just to check if we were given valid json by Redis
			json, err := simplejson.NewJson(value)
			if err != nil {
				errorMessage := "Error while unmarshalling json from Redis: " + err.Error()
				fmt.Println(LogTag, errorMessage)
				serveJSONError(w, errorMessage)
				return 0, nil
			}
			// Check what type of encoding we need to serve
			if strings.Contains(r.Header.Get("Accept-Encoding"), "xml") {
				// Create XML
				jsonMap, err := json.Map()
				if err != nil {
					errorMessage := "Error while converting json to Golang map: " + err.Error()
					fmt.Println(LogTag, errorMessage)
					serveXMLError(w, errorMessage)
					return 0, nil
				}
				xmlMap := mxj.Map(jsonMap)
				xmlValue, err := xmlMap.Xml()
				if err != nil {
					errorMessage := "Error while unmarshalling xml: " + err.Error()
					fmt.Println(LogTag, errorMessage)
					serveXMLError(w, errorMessage)
					return 0, nil
				}
				fmt.Println("XML Payload:", xmlValue)
				serveXML(w, xmlValue)
			} else {
				serveJSON(w, []byte(value))
			}
		}
		return http.StatusOK, nil
	} else {
		return redis.Next.ServeHTTP(w, r)
	}
}

func serveJSONError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("{\"error\":\"" + errorMessage + "\"}"))
}

func serveJSON(w http.ResponseWriter, jsonPayload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPayload)
}

func serveXMLError(w http.ResponseWriter, errorMessage string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\"?><error>" + errorMessage + "</error>"))
}

func serveXML(w http.ResponseWriter, xmlPayload []byte) {
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xmlPayload)
}
