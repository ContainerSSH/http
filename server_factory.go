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
	onReady func(string),
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
		onReady:   onReady,
	}, nil
}

func createServerTLSConfig(config ServerConfiguration) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		MinVersion:               config.TLSVersion.getTLSVersion(),
		CurvePreferences:         config.ECDHCurves.getList(),
		PreferServerCipherSuites: true,
		CipherSuites:             config.CipherSuites.getList(),
	}

	tlsConfig.Certificates = []tls.Certificate{*config.cert}

	if config.clientCAPool != nil {
		tlsConfig.ClientCAs = config.clientCAPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}
