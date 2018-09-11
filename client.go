package gohttpproxy

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Client  *http.Client
	Timeout time.Duration
}

// Proxy do Client work, request real server to get response, need a complete new
// header, to set Proxy-Connection by golang http lib
func (c *Client) Request(r *http.Request) (*http.Response, error) {

	// path := getURI(r)
	req, err := http.NewRequest(r.Method, r.RequestURI, r.Body)

	if err != nil {
		log.Println("form new request failed")
		return nil, errors.New("Form new request failed")
	}

	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	for _, v := range r.Cookies() {
		req.Header.Add("Cookie", v.Raw)
	}

	//log.Println("send:", req)
	resp, err := c.Client.Do(req)
	log.Println("resp received")
	return resp, err
}

func getURI(r *http.Request) string {
	uriInfo, err := url.ParseRequestURI(r.RequestURI)
	if err != nil {
		return ""
	}
	// log.Println("getURI:", uriInfo.Scheme)
	// log.Println("getURI:", uriInfo.Opaque)
	// log.Println("getURI:", uriInfo.RawQuery)
	// log.Println("getURI:", uriInfo.Host)
	log.Println("getURI:", uriInfo.Path)
	// log.Println("getURI:", uriInfo.RawPath)
	// log.Println("getURI:", uriInfo.RawQuery)
	// log.Println("getURI:", uriInfo.Fragment)
	return uriInfo.Path
}
