package heroku

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
)

// NewTLSConfig constructs a *tls.Config from the given certificates that is
// appropriate for connecting to a Heroku Kafka broker.
func NewTLSConfig(trustedCert string, clientCert string, clientKey string) (*tls.Config, error) {
	roots := x509.NewCertPool()

	// Build a certificate authority pool for verification
	ok := roots.AppendCertsFromPEM([]byte(trustedCert))
	if !ok {
		return nil, errors.New("Invalid Trusted Cert")
	}

	// Parse the client certificate key pair
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,

		// Disable normal certificate verification because Heroku's certificates do
		// not match their hostnames. Even when enabled, the customer
		// VerifyPeerCertificate will be called.
		InsecureSkipVerify: true,

		// VerifyPeerCertificate will check that the certificate was signed by the
		// trusted certificate Heroku sets in the environment. The second parameter
		// will always be empty since InsecureSkipVerify is true.
		VerifyPeerCertificate: func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
			opts := x509.VerifyOptions{Roots: roots}

			// Heroku is only giving one rawCert as of 2017-09-19. If that changes
			// then this code may no longer be valid.
			for _, raw := range rawCerts {
				cert, err := x509.ParseCertificate(raw)
				if err != nil {
					return err
				}

				// If any of the raw certificates fail verification then we fail the
				// entire process.
				_, err = cert.Verify(opts)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cfg, nil
}
