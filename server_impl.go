package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	goHttp "net/http"
	"sync"
)

type server struct {
	lock      *sync.Mutex
	handler   goHttp.Handler
	config    ServerConfiguration
	tlsConfig *tls.Config
	srv       *goHttp.Server
	done      chan bool
	goLogger  io.Writer
	onReady   func()
}

func (s *server) Run() error {
	s.lock.Lock()
	if s.srv != nil {
		return fmt.Errorf("server is already running")
	}
	s.srv = &goHttp.Server{
		Addr:      s.config.Listen,
		Handler:   s.handler,
		TLSConfig: s.tlsConfig,
		ErrorLog:  log.New(s.goLogger, "", 0),
	}
	defer func() {
		s.lock.Lock()
		s.srv = nil
		s.lock.Unlock()
		s.done <- true
	}()
	var err error

	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		s.lock.Unlock()
		return err
	}
	defer func() { _ = ln.Close() }()
	if s.onReady != nil {
		s.onReady()
	}
	s.lock.Unlock()
	if s.srv.TLSConfig != nil {
		err = s.srv.ServeTLS(ln, "", "")
	} else {
		err = s.srv.Serve(ln)
	}
	if err != nil && !errors.Is(err, goHttp.ErrServerClosed) {
		return err
	}
	return nil
}

func (s *server) Shutdown(shutdownContext context.Context) {
	s.lock.Lock()
	if s.srv == nil {
		s.lock.Unlock()
		return
	}
	srv := s.srv
	done := s.done
	s.lock.Unlock()
	// Ignore error because we don't care about shutdown context violations, we wait anyway
	_ = srv.Shutdown(shutdownContext)
	<-done
}
