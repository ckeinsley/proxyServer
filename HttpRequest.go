package main

import (
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"
)

// HTTPRequest represents an arbitrary HTTP request
type HTTPRequest struct {
	Version    string // Always HTTP/1.0
	Method     string // Type of request [GET, POST, ...]
	Host       string
	Port       string
	Route      string // Requested page
	Headers    []string
	Connection string // Always Connection: close
	Hash       string
}

// CreateHTTPRequest creates an HTTPRequest object out of the string sent to the proxy
// example of input below
func CreateHTTPRequest(connectionRequest string) HTTPRequest {
	spaceSplitRequest := strings.Split(connectionRequest, " ")
	method := spaceSplitRequest[0]
	noHTTP := strings.Replace(spaceSplitRequest[1], "http://", "", 1)

	slicedOnSlash := strings.SplitN(noHTTP, "/", 2)
	route := "/"
	if len(slicedOnSlash) > 1 {
		route += slicedOnSlash[1]
	}
	re := regexp.MustCompile(`\r?\n`)
	route = re.ReplaceAllString(route, "")

	hostWithPort := slicedOnSlash[0]
	hostPortSlice := strings.Split(hostWithPort, ":")
	host := hostPortSlice[0]
	port := "80"
	if len(hostPortSlice) > 1 {
		port = hostPortSlice[1]
	}

	return HTTPRequest{
		Method:     method,
		Host:       host,
		Route:      route,
		Version:    "HTTP/1.0",
		Port:       port,
		Headers:    parseHeaders(connectionRequest),
		Connection: "Connection: close",
		Hash:       hash(host + port + route + method),
	}
}

func parseHeaders(connectionRequest string) []string {
	var headers []string

	tokens := strings.Split(connectionRequest, "\n")
	for index, token := range tokens {
		if index < 2 {
			continue
		}
		if strings.Contains(token, "Connection:") {
			continue
		}
		headers = append(headers, token)
	}

	return headers
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return strconv.Itoa(int(h.Sum32()))
}

/* Example Requests coming into proxy */

// GET http://www.columbia.edu/~fdc/sample.html HTTP/1.0
// Host: www.columbia.edu
// User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0
// Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
// Accept-Language: en-US,en;q=0.5
// Accept-Encoding: gzip, deflate
// Connection: keep-alive
// Upgrade-Insecure-Requests: 1
// Pragma: no-cache

// GET http://www.columbia.edu/favicon.ico HTTP/1.0
// Host: www.columbia.edu
// User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0
// Accept: image/webp,*/*
// Accept-Language: en-US,en;q=0.5
// Accept-Encoding: gzip, deflate
// Connection: keep-alive
// Pragma: no-cache

// GET http://detectportal.firefox.com/success.txt HTTP/1.0
// Host: detectportal.firefox.com
// User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0
// Accept: */*
// Accept-Language: en-US,en;q=0.5
// Accept-Encoding: gzip, deflate
// Cache-Control: no-cache
// Pragma: no-cache
// Connection: keep-alive
