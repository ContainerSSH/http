package http

import (
	"crypto/tls"
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

	return &client{
		config:    config,
		logger:    logger.WithLabel("endpoint", config.URL),
		tlsConfig: tlsConfig,
	}, nil
}

// createTLSConfig creates a TLS config. Should only be called after config.Validate().
func createTLSConfig(config ClientConfiguration) (*tls.Config, error) {
	if !strings.HasPrefix(config.URL, "https://") {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion:       config.TLSVersion.getTLSVersion(),
		CurvePreferences: config.ECDHCurves.getList(),
		CipherSuites:     config.CipherSuites.getList(),
	}
	if config.caCertPool != nil {
		tlsConfig.RootCAs = config.caCertPool
	}
	if config.cert != nil {
		tlsConfig.Certificates = []tls.Certificate{*config.cert}
	}
	return tlsConfig, nil
}
