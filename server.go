package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"

	"github.com/containerssh/service"
)

// ServerConfiguration is a structure to configure the simple HTTP server by.
//goland:noinspection GoVetStructTag
type ServerConfiguration struct {
	// Listen contains the IP and port to listen on.
	Listen string `json:"listen" yaml:"listen" default:"0.0.0.0:8080"`
	// Key contains either a file name to a private key, or the private key itself in PEM format to use as a server key.
	Key string `json:"key" yaml:"key"`
	// Cert contains either a file to a certificate, or the certificate itself in PEM format to use as a server
	// certificate.
	Cert string `json:"cert" yaml:"cert"`
	// ClientCACert contains either a file or a certificate in PEM format to verify the connecting clients by.
	ClientCACert string `json:"clientcacert" yaml:"clientcacert"`

	// cert is for internal use only. It contains the key and certificate after Validate.
	cert *tls.Certificate `json:"-" yaml:"-"`
	// clientCAPool is for internal use only. It contains the client CA pool after Validate.
	clientCAPool *x509.CertPool `json:"-" yaml:"-"`
}

// Validate validates the server configuration.
func (config *ServerConfiguration) Validate() error {
	if config.Listen == "" {
		return fmt.Errorf("no listen address provided")
	}
	if _, _, err := net.SplitHostPort(config.Listen); err != nil {
		return fmt.Errorf("invalid listen address provided (%w)", err)
	}
	if config.Cert != "" && config.Key == "" {
		return fmt.Errorf("certificate provided without a key")
	}
	if config.Cert == "" && config.Key != "" {
		return fmt.Errorf("key provided without certificate")
	}

	if config.Cert != "" && config.Key != "" {
		pemCert, err := loadPem(config.Cert)
		if err != nil {
			return fmt.Errorf("failed to load certificate (%w)", err)
		}
		pemKey, err := loadPem(config.Key)
		if err != nil {
			return fmt.Errorf("failed to load key (%w)", err)
		}
		cert, err := tls.X509KeyPair(pemCert, pemKey)
		if err != nil {
			return fmt.Errorf("failed to load key/certificate (%w)", err)
		}
		config.cert = &cert
	}

	if config.ClientCACert != "" {
		clientCaCert, err := loadPem(config.ClientCACert)
		if err != nil {
			return fmt.Errorf("failed to load client CA certificate (%w)", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(clientCaCert) {
			return fmt.Errorf("failed to load client CA certificate")
		}
		config.clientCAPool = caCertPool
	}

	return nil
}

// Server is an interface that specifies the minimum requirements for the server.
type Server interface {
	service.Service
}
