package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	goHttp "net/http"
	"sync"

	"github.com/containerssh/service"
)

type server struct {
	name      string
	lock      *sync.Mutex
	handler   goHttp.Handler
	config    ServerConfiguration
	tlsConfig *tls.Config
	srv       *goHttp.Server
	goLogger  io.Writer
}

func (s *server) String() string {
	return s.name
}

func (s *server) RunWithLifecycle(lifecycle service.Lifecycle) error {
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
	}()
	var err error

	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		s.lock.Unlock()
		return err
	}
	defer func() { _ = ln.Close() }()
	lifecycle.Running()
	s.lock.Unlock()
	serverFinished := make(chan struct{}, 1)
	go func() {
		select {
		case <-lifecycle.Context().Done():
			s.lock.Lock()
			if s.srv == nil {
				s.lock.Unlock()
				return
			}
			srv := s.srv
			s.lock.Unlock()
			_ = srv.Shutdown(lifecycle.ShutdownContext())
		case <-serverFinished:
		}
	}()
	if s.srv.TLSConfig != nil {
		err = s.srv.ServeTLS(ln, "", "")
	} else {
		err = s.srv.Serve(ln)
	}
	serverFinished <- struct{}{}
	if err != nil && !errors.Is(err, goHttp.ErrServerClosed) {
		return err
	}
	return nil
}
