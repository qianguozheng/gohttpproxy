package gohttpproxy

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	stdLog "log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/labstack/gommon/color"
	"github.com/labstack/gommon/log"
)

/// GoHTTPProxy
/// Server listen to socket, and receive http requests, parse http header, then relay to
/// client to form new http request, send to server.

/// DNS client to parse host to ip address

/// Client connect to target server and get response from server, relay back to Server.

type Proxy struct {
	stdLogger   *stdLog.Logger
	colorer     *color.Color
	Server      *http.Server
	TLSServer   *http.Server
	Listener    net.Listener
	TLSListener net.Listener
	Debug       bool
	Mutex       sync.RWMutex
	Logger      Logger
}

// New creates an instance of Proxy
func New() (p *Proxy) {
	p = &Proxy{
		Server:    new(http.Server),
		TLSServer: new(http.Server),
		Logger:    log.New("Proxy"),
		colorer:   color.New(),
	}

	p.Server.Handler = p
	p.TLSServer.Handler = p

	p.Logger.SetLevel(log.OFF)
	p.stdLogger = stdLog.New(p.Logger.Output(), p.Logger.Prefix()+": ", 0)
	return
}

// Start starts an HTTP server.
func (p *Proxy) Start(address string) error {
	p.Server.Addr = address
	return p.StartServer(p.Server)
}

// StartServer starts a custom http server.
func (p *Proxy) StartServer(s *http.Server) (err error) {
	// Setup
	p.colorer.SetOutput(p.Logger.Output())
	s.Handler = p
	s.ErrorLog = p.stdLogger

	if s.TLSConfig == nil {
		if p.Listener == nil {
			p.Listener, err = newListener(s.Addr)
			if err != nil {
				return err
			}
		}
		p.colorer.Printf("=> http server started on %s\n", p.colorer.Green(p.Listener.Addr()))
		return s.Serve(p.Listener)
	}
	if p.TLSListener == nil {
		l, err := newListener(s.Addr)
		if err != nil {
			return err
		}
		p.TLSListener = tls.NewListener(l, s.TLSConfig)
	}
	p.colorer.Printf("⇛ https server started on %s\n", p.colorer.Green(p.TLSListener.Addr()))
	return s.Serve(p.TLSListener)

}

// ServeHTTP implements `http.Handler` interface, which serves HTTP requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Acquire lock
	p.Mutex.RLock()
	defer p.Mutex.RUnlock()

	//TODO: Do Proxy Things
	fmt.Println("ServeHTTP Interface")
	c := &Client{Client: &http.Client{}}
	resp, err := c.Request(r)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	for k, v := range resp.Header {
		for _, vv := range v {
			w.Header().Add(k, vv)
		}
	}
	w.WriteHeader(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Response failed ", err.Error())
	}
	w.Write(body)
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}
