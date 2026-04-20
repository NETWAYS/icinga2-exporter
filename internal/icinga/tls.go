package icinga

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type closeIdler interface {
	CloseIdleConnections()
}

// cloneRequest returns a clone of the provided *http.Request
func cloneRequest(r *http.Request) *http.Request {
	// Shallow copy of the struct.
	r2 := new(http.Request)
	*r2 = *r
	// Deep copy of the Header.
	r2.Header = make(http.Header)
	maps.Copy(r2.Header, r.Header)

	return r2
}

// readCAFile reads the CA cert file from disk.
func readCAFile(f string) ([]byte, error) {
	d, err := os.ReadFile(filepath.Clean(f))

	if err != nil {
		return nil, fmt.Errorf("unable to load CA cert %s: %w", f, err)
	}

	return d, nil
}

// getClientCertificate reads the pair of client cert and key from disk and returns a tls.Certificate.
func (c *TLSConfig) getClientCertificate(_ *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	certData, keyData, err := readCertAndKey(c.CertFile, c.KeyFile)

	if err != nil {
		return nil, fmt.Errorf("unable to read client cert (%s) & key (%s): %w", c.CertFile, c.KeyFile, err)
	}

	cert, err := tls.X509KeyPair(certData, keyData)

	if err != nil {
		return nil, fmt.Errorf("unable to use client cert (%s) & key (%s): %w", c.CertFile, c.KeyFile, err)
	}

	return &cert, nil
}

// readCertAndKey reads the cert and key files from the disk.
func readCertAndKey(certFile, keyFile string) ([]byte, []byte, error) {
	certData, err := os.ReadFile(filepath.Clean(certFile))

	if err != nil {
		return nil, nil, err
	}

	keyData, err := os.ReadFile(filepath.Clean(keyFile))
	if err != nil {
		return nil, nil, err
	}

	return certData, keyData, nil
}

// updateRootCA parses the given byte slice as a series of PEM encoded certificates and updates tls.Config.RootCAs.
func updateRootCA(cfg *tls.Config, b []byte) bool {
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(b) {
		return false
	}

	cfg.RootCAs = caCertPool

	return true
}

// TLSConfig configures the options for TLS connections.
type TLSConfig struct {
	CAFile             string
	CertFile           string
	KeyFile            string
	ServerName         string
	InsecureSkipVerify bool
}

// NewTLSConfig creates a new tls.Config from the given TLSConfig.
func newTLSConfig(cfg *TLSConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		// nolint: gosec
		InsecureSkipVerify: cfg.InsecureSkipVerify,
	}

	// If a CA cert is provided
	if len(cfg.CAFile) > 0 {
		b, err := readCAFile(cfg.CAFile)

		if err != nil {
			return nil, err
		}

		if !updateRootCA(tlsConfig, b) {
			return nil, fmt.Errorf("unable to use CA cert %s", cfg.CAFile)
		}
	}

	if len(cfg.ServerName) > 0 {
		tlsConfig.ServerName = cfg.ServerName
	}

	// If a client cert & key is provided then configure TLS config accordingly.
	// nolint: gocritic
	if len(cfg.CertFile) > 0 && len(cfg.KeyFile) == 0 {
		return nil, fmt.Errorf("client cert file %q specified without client key file", cfg.CertFile)
	} else if len(cfg.KeyFile) > 0 && len(cfg.CertFile) == 0 {
		return nil, fmt.Errorf("client key file %q specified without client cert file", cfg.KeyFile)
	} else if len(cfg.CertFile) > 0 && len(cfg.KeyFile) > 0 {
		// Verify that client cert and key are valid.
		if _, err := cfg.getClientCertificate(nil); err != nil {
			return nil, err
		}

		tlsConfig.GetClientCertificate = cfg.getClientCertificate
	}

	return tlsConfig, nil
}

type authorizationCredentialsRoundTripper struct {
	authType        string
	authCredentials string
	rt              http.RoundTripper
}

// NewAuthorizationCredentialsRoundTripper adds the provided credentials to a
// request unless the authorization header has already been set.
func NewAuthorizationCredentialsRoundTripper(authType, authCredentials string, rt http.RoundTripper) http.RoundTripper {
	return &authorizationCredentialsRoundTripper{authType, authCredentials, rt}
}

func (rt *authorizationCredentialsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) == 0 {
		req = cloneRequest(req)
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", rt.authType, rt.authCredentials))
	}

	return rt.rt.RoundTrip(req)
}

func (rt *authorizationCredentialsRoundTripper) CloseIdleConnections() {
	if ci, ok := rt.rt.(closeIdler); ok {
		ci.CloseIdleConnections()
	}
}

type basicAuthRoundTripper struct {
	username string
	password string
	rt       http.RoundTripper
}

// NewBasicAuthRoundTripper adds the provided basic auth credentials to a request
func newBasicAuthRoundTripper(username, password string, rt http.RoundTripper) http.RoundTripper {
	return &basicAuthRoundTripper{username, password, rt}
}

func (rt *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(req.Header.Get("Authorization")) != 0 {
		return rt.rt.RoundTrip(req)
	}

	req = cloneRequest(req)

	req.SetBasicAuth(rt.username, strings.TrimSpace(rt.password))

	return rt.rt.RoundTrip(req)
}

func (rt *basicAuthRoundTripper) CloseIdleConnections() {
	if ci, ok := rt.rt.(closeIdler); ok {
		ci.CloseIdleConnections()
	}
}
