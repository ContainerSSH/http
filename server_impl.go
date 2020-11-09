package http

import (
	"context"
	"crypto/tls"
	"fmt"
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
}

func (s *server) Run() error {
	s.lock.Lock()
	if s.srv != nil {
		return fmt.Errorf("server is already running")
	}
	s.srv = &goHttp.Server{
		Addr:      s.config.Listen,
		Handler:   s.handler,
		TLSConfig: nil,
	}
	defer func() {
		s.lock.Lock()
		s.srv = nil
		s.lock.Unlock()
		s.done <- true
	}()
	s.lock.Unlock()
	var err error
	if s.srv.TLSConfig != nil {
		err = s.srv.ListenAndServeTLS("", "")
	} else {
		err = s.srv.ListenAndServe()
	}
	if err != nil && err != goHttp.ErrServerClosed {
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
