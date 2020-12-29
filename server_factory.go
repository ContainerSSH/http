package http

import (
	"crypto/tls"
	goHttp "net/http"
	"sync"

	"github.com/containerssh/log"
)

// NewServer creates a new HTTP server with the given configuration and calling the provided handler.
func NewServer(
	name string,
	config ServerConfiguration,
	handler goHttp.Handler,
	logger log.Logger,
) (Server, error) {
	if handler == nil {
		panic("BUG: no handler provided to http.NewServer")
	}
	if logger == nil {
		panic("BUG: no logger provided to http.NewServer")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	var tlsConfig *tls.Config
	if config.cert != nil {
		var err error
		tlsConfig, err = createServerTLSConfig(config)
		if err != nil {
			return nil, err
		}
	}

	return &server{
		name:      name,
		lock:      &sync.Mutex{},
		handler:   handler,
		config:    config,
		tlsConfig: tlsConfig,
		srv:       nil,
		goLogger:  log.NewGoLogWriter(logger),
	}, nil
}

func createServerTLSConfig(config ServerConfiguration) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS13,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	tlsConfig.Certificates = []tls.Certificate{*config.cert}

	if config.clientCAPool != nil {
		tlsConfig.ClientCAs = config.clientCAPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}
