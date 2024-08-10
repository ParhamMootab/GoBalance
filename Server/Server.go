package Server

import (
	"net/http/httputil"
	"net/url"
)


type Server struct {
	URL          *url.URL
	ReverseProxy *httputil.ReverseProxy
	Healthy      bool
	Weight       int
}