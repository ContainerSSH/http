package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	goHttp "net/http"
	"sync"

	"github.com/containerssh/log"
)

// NewServer creates a new HTTP server with the given configuration and calling the provided handler.
func NewServer(
	config ServerConfiguration,
	handler goHttp.Handler,
	onReady func(),
	logger log.Logger,
) (Server, error) {
	if handler == nil {
		return nil, fmt.Errorf("handler cannot be nil")
	}

	var tlsConfig *tls.Config
	if config.Cert != "" && config.Key != "" {
		var err error
		tlsConfig, err = createServerTlsConfig(config)
		if err != nil {
			return nil, err
		}
	} else {
		if config.Cert != "" {
			return nil, fmt.Errorf("server certificate provided, but no private key")
		}
		if config.Key != "" {
			return nil, fmt.Errorf("server privaet key provided, but no certificate")
		}
		if config.ClientCaCert != "" {
			return nil, fmt.Errorf("client CA certificate is set, but no server certificate or private key provided")
		}
	}

	return &server{
		lock:      &sync.Mutex{},
		handler:   handler,
		config:    config,
		tlsConfig: tlsConfig,
		srv:       nil,
		done:      make(chan bool, 1),
		goLogger:  log.NewGoLogWriter(logger),
		onReady:   onReady,
	}, nil
}

func createServerTlsConfig(config ServerConfiguration) (*tls.Config, error) {
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

	clientCert, err := loadPem(config.Cert)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate (%w)", err)
	}
	clientKey, err := loadPem(config.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate (%w)", err)
	}
	cert, err := tls.X509KeyPair(clientCert, clientKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate or key (%w)", err)
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	if config.ClientCaCert != "" {
		clientCaCert, err := loadPem(config.ClientCaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate (%w)", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(clientCaCert)
		tlsConfig.ClientCAs = caCertPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}
	return tlsConfig, nil
}
