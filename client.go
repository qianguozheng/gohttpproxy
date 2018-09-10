package gohttpproxy

import (
	"errors"
	"log"
	"net/http"
)

type Client struct {
	Client *http.Client
}

// Proxy do Client work, request real server to get response, need a complete new
// header, to set Proxy-Connection by golang http lib
func (c *Client) Request(r *http.Request) (*http.Response, error) {
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
	//log.Println("send:", req)
	resp, err := c.Client.Do(req)
	log.Println("resp received")
	return resp, err
}
