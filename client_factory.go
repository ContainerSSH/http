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

func NewClient(
	config ClientConfiguration,
	logger log.Logger,
) (Client, error) {
	tlsConfig, err := createTlsConfig(config, logger)
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{TLSClientConfig: tlsConfig}

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

func createTlsConfig(config ClientConfiguration, logger log.Logger) (*tls.Config, error) {
	if logger == nil {
		return nil, fmt.Errorf("no logger provided")
	}

	if config.Url == "" {
		return nil, fmt.Errorf("no URL provided")
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
		cert, err := tls.LoadX509KeyPair(string(clientCert), string(clientKey))
		if err != nil {
			logger.Criticale(err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}
