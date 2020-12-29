package http

import (
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/containerssh/log"
)

// NewClient creates a new HTTP client with the given configuration.
func NewClient(
	config ClientConfiguration,
	logger log.Logger,
) (Client, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}
	if logger == nil {
		panic("BUG: no logger provided for http.NewClient")
	}

	tlsConfig, err := createTLSConfig(config)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &client{
		config:     config,
		logger:     logger,
		tlsConfig:  tlsConfig,
		httpClient: httpClient,
	}, nil
}

// createTLSConfig creates a TLS config. Should only be called after config.Validate().
func createTLSConfig(config ClientConfiguration) (*tls.Config, error) {
	if !strings.HasPrefix(config.URL, "https://") {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS13,
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP521, tls.CurveP384, tls.CurveP256},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}
	if config.caCertPool != nil {
		tlsConfig.RootCAs = config.caCertPool
	}
	if config.cert != nil {
		tlsConfig.Certificates = []tls.Certificate{*config.cert}
	}
	return tlsConfig, nil
}
