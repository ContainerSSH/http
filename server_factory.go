package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	goHttp "net/http"
	"strings"
	"sync"

	"github.com/containerssh/log"
)

// NewServer creates a new HTTP server with the given configuration and calling the provided handler.
func NewServer(
	config ServerConfiguration,
	handler goHttp.Handler,
	logger log.Logger,
) (Server, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	}

	tlsConfig := &tls.Config{}
	if strings.TrimSpace(config.CaCert) != "" {
		caCert, err := loadPem(config.CaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate (%w)", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	if config.Cert != "" && config.Key != "" {
		clientCert, err := loadPem(config.Cert)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate (%w)", err)
		}
		clientKey, err := loadPem(config.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate (%w)", err)
		}
		cert, err := tls.LoadX509KeyPair(string(clientCert), string(clientKey))
		if err != nil {
			logger.Criticale(err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if config.ClientCaCert != "" {
		caCert, err := loadPem(config.CaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate (%w)", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return &server{
		lock:      &sync.Mutex{},
		handler:   handler,
		config:    config,
		tlsConfig: tlsConfig,
		srv:       nil,
		done:      make(chan bool, 1),
	}, nil
}
