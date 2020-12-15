package http_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"testing"
	"time"

	"github.com/containerssh/log"
	"github.com/containerssh/service"
	"github.com/stretchr/testify/assert"

	"github.com/containerssh/http"
)

type Request struct {
	Message string `json:"Message"`
}

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"Message"`
}

type handler struct {
}

func (s *handler) OnRequest(request http.ServerRequest, response http.ServerResponse) error {
	req := Request{}
	if err := request.Decode(&req); err != nil {
		return err
	}
	if req.Message == "Hi" {
		response.SetBody(&Response{
			Error:   false,
			Message: "Hello world!",
		})
	} else {
		response.SetStatus(400)
		response.SetBody(&Response{
			Error:   true,
			Message: "Be nice and greet me!",
		})
	}
	return nil
}

func TestUnencrypted(t *testing.T) {
	clientConfig := http.ClientConfiguration{
		URL:     "http://127.0.0.1:8080/",
		Timeout: 2 * time.Second,
	}
	serverConfig := http.ServerConfiguration{
		Listen: "127.0.0.1:8080",
	}

	message := "Hi"

	response, responseStatus, err := runRequest(clientConfig, serverConfig, message)
	if err != nil {
		assert.Fail(t, "failed to run request", err)
		return
	}
	assert.Equal(t, 200, responseStatus)
	assert.Equal(t, false, response.Error)
	assert.Equal(t, "Hello world!", response.Message)
}

func TestUnencryptedFailure(t *testing.T) {
	clientConfig := http.ClientConfiguration{
		URL:     "http://127.0.0.1:8080/",
		Timeout: 2 * time.Second,
	}
	serverConfig := http.ServerConfiguration{
		Listen: "127.0.0.1:8080",
	}

	message := "Hm..."

	response, responseStatus, err := runRequest(clientConfig, serverConfig, message)
	if err != nil {
		assert.Fail(t, "failed to run request", err)
		return
	}
	assert.Equal(t, 400, responseStatus)
	assert.Equal(t, true, response.Error)
	assert.Equal(t, "Be nice and greet me!", response.Message)
}

func TestEncrypted(t *testing.T) {
	caPrivKey, caCert, caCertBytes, err := createCA()
	if err != nil {
		assert.Fail(t, "failed to create CA", err)
		return
	}
	serverPrivKey, serverCert, err := createSignedCert(
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		caPrivKey,
		caCert,
	)
	if err != nil {
		assert.Fail(t, "failed to create server cert", err)
		return
	}

	clientConfig := http.ClientConfiguration{
		URL:     "https://127.0.0.1:8080/",
		Timeout: 2 * time.Second,
		CaCert:  string(caCertBytes),
	}
	serverConfig := http.ServerConfiguration{
		Listen: "127.0.0.1:8080",
		Key:    string(serverPrivKey),
		Cert:   string(serverCert),
	}

	message := "Hi"

	response, responseStatus, err := runRequest(clientConfig, serverConfig, message)
	if err != nil {
		assert.Fail(t, "failed to run request", err)
		return
	}
	assert.Equal(t, 200, responseStatus)
	assert.Equal(t, false, response.Error)
	assert.Equal(t, "Hello world!", response.Message)
}

func TestMutuallyAuthenticated(t *testing.T) {
	caPrivKey, caCert, caCertBytes, err := createCA()
	if err != nil {
		assert.Fail(t, "failed to create CA", err)
		return
	}
	serverPrivKey, serverCert, err := createSignedCert(
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		caPrivKey,
		caCert,
	)
	if err != nil {
		assert.Fail(t, "failed to create server cert", err)
		return
	}

	clientCaPriv, clientCaCert, clientCaCertBytes, err := createCA()
	if err != nil {
		assert.Fail(t, "failed to create client CA", err)
		return
	}
	clientPrivKey, clientCert, err := createSignedCert(
		[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		clientCaPriv,
		clientCaCert,
	)
	if err != nil {
		assert.Fail(t, "failed to create server cert", err)
		return
	}

	clientConfig := http.ClientConfiguration{
		URL:        "https://127.0.0.1:8080/",
		CaCert:     string(caCertBytes),
		Timeout:    2 * time.Second,
		ClientCert: string(clientCert),
		ClientKey:  string(clientPrivKey),
	}
	serverConfig := http.ServerConfiguration{
		Listen:       "127.0.0.1:8080",
		Key:          string(serverPrivKey),
		Cert:         string(serverCert),
		ClientCaCert: string(clientCaCertBytes),
	}

	message := "Hi"

	response, responseStatus, err := runRequest(clientConfig, serverConfig, message)
	if err != nil {
		assert.Fail(t, "failed to run request", err)
		return
	}
	assert.Equal(t, 200, responseStatus)
	assert.Equal(t, false, response.Error)
	assert.Equal(t, "Hello world!", response.Message)
}

func TestMutuallyAuthenticatedFailure(t *testing.T) {
	caPrivKey, caCert, caCertBytes, err := createCA()
	if err != nil {
		assert.Fail(t, "failed to create CA", err)
		return
	}
	serverPrivKey, serverCert, err := createSignedCert(
		[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		caPrivKey,
		caCert,
	)
	if err != nil {
		assert.Fail(t, "failed to create server cert", err)
		return
	}

	clientCaPriv, clientCaCert, _, err := createCA()
	if err != nil {
		assert.Fail(t, "failed to create client CA", err)
		return
	}
	clientPrivKey, clientCert, err := createSignedCert(
		[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		clientCaPriv,
		clientCaCert,
	)
	if err != nil {
		assert.Fail(t, "failed to create server cert", err)
		return
	}

	clientConfig := http.ClientConfiguration{
		URL:        "https://127.0.0.1:8080/",
		CaCert:     string(caCertBytes),
		Timeout:    2 * time.Second,
		ClientCert: string(clientCert),
		ClientKey:  string(clientPrivKey),
	}
	serverConfig := http.ServerConfiguration{
		Listen: "127.0.0.1:8080",
		Key:    string(serverPrivKey),
		Cert:   string(serverCert),
		//Pass wrong client CA cert to test failure
		ClientCaCert: string(caCertBytes),
	}

	message := "Hi"

	if _, _, err = runRequest(clientConfig, serverConfig, message); err == nil {
		assert.Fail(t, "Client request with invalid CA verification did not fail.")
		return
	}
	println(err)
}

func createCA() (*rsa.PrivateKey, *x509.Certificate, []byte, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"ACME, Inc"},
			Country:      []string{"US"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create private key (%w)", err)
	}
	caCert, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to create CA certificate (%w)", err)
	}
	caPEM := new(bytes.Buffer)
	if err := pem.Encode(
		caPEM,
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caCert,
		},
	); err != nil {
		return nil, nil, nil, fmt.Errorf("failed to encode CA cert (%w)", err)
	}
	return caPrivateKey, ca, caPEM.Bytes(), nil
}

func createSignedCert(
	usage []x509.ExtKeyUsage,
	caPrivateKey *rsa.PrivateKey,
	caCertificate *x509.Certificate,
) ([]byte, []byte, error) {
	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"ACME, Inc"},
			Country:      []string{"US"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1)},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(0, 0, 1),
		SubjectKeyId: []byte{1},
		ExtKeyUsage:  usage,
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		cert,
		caCertificate,
		&certPrivKey.PublicKey,
		caPrivateKey,
	)
	if err != nil {
		return nil, nil, err
	}
	certPrivKeyPEM := new(bytes.Buffer)
	if err := pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	}); err != nil {
		return nil, nil, err
	}
	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM,
		&pem.Block{Type: "CERTIFICATE", Bytes: certBytes},
	); err != nil {
		return nil, nil, err
	}
	return certPrivKeyPEM.Bytes(), certPEM.Bytes(), nil
}

func runRequest(
	clientConfig http.ClientConfiguration,
	serverConfig http.ServerConfiguration,
	message string,
) (Response, int, error) {
	response := Response{}
	logger, err := log.New(
		log.Config{
			Level:  log.LevelDebug,
			Format: log.FormatText,
		},
		"http",
		os.Stdout,
	)
	if err != nil {
		return response, 0, err
	}
	client, err := http.NewClient(clientConfig, logger)
	if err != nil {
		return response, 0, err
	}

	ready := make(chan bool, 1)
	server, err := http.NewServer(
		"HTTP",
		serverConfig,
		http.NewServerHandler(&handler{}, logger),
		logger,
	)
	if err != nil {
		return response, 0, err
	}
	lifecycle := service.NewLifecycle(server)
	lifecycle.OnRunning(func(s service.Service, l service.Lifecycle) {
		ready <- true
	})

	errorChannel := make(chan error, 2)
	responseStatus := 0
	go func() {
		if err := lifecycle.Run(); err != nil {
			errorChannel <- err
		}
		close(errorChannel)
	}()
	<-ready
	if responseStatus, err = client.Post(
		context.Background(),
		"",
		&Request{Message: message},
		&response,
	); err != nil {
		errorChannel <- err
	}
	lifecycle.Stop(context.Background())
	if err, ok := <-errorChannel; ok {
		return response, 0, err
	}
	return response, responseStatus, nil
}
