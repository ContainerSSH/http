package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/containerssh/log"
)

// NewClient creates a new HTTP client with the given configuration.
func NewClient(
	config ClientConfiguration,
	logger log.Logger,
) (Client, error) {
	tlsConfig, err := createTLSConfig(config, logger)
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

func createTLSConfig(config ClientConfiguration, logger log.Logger) (*tls.Config, error) {
	if logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}

	if config.Url == "" {
		return nil, fmt.Errorf("no URL provided")
	}

	if !strings.HasPrefix(config.Url, "https://") {
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
	if strings.TrimSpace(config.CaCert) != "" {
		caCert, err := loadPem(config.CaCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate (%w)", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	} else if runtime.GOOS == "windows" && strings.HasPrefix(config.Url, "https://") {
		//Remove if https://github.com/golang/go/issues/16736 gets fixed
		return nil, fmt.Errorf(
			"due to a bug (#16736) in Golang on Windows CA certificates have to be explicitly " +
				"provided for https:// authentication server URLs",
		)
	}

	if config.ClientCert != "" && config.ClientKey != "" {
		clientCert, err := loadPem(config.ClientCert)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate (%w)", err)
		}
		clientKey, err := loadPem(config.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate (%w)", err)
		}
		cert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate or key (%w)", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}
