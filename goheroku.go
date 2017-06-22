package goheroku

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"net/url"
	"os"
	"strings"
)

// To specify the topic or consumer group, the Kafka prefix needs
// to appended to it. This function makes it possible without having
// access to the app config.
func AppendPrefixTo(str string) string {
	prefix := os.Getenv("KAFKA_PREFIX")
	if prefix != "" {
		str = strings.Join([]string{prefix, str}, "")
	}
	return str
}

// Create the TLS context, using the key and certificates provided.
func createTLSConfig() (*tls.Config, error) {
	trustedCert, ok := os.LookupEnv("KAFKA_TRUSTED_CERT")
	if !ok {
		return nil, errors.New("Kafka Trusted Certificate not found!")
	}

	clientCertKey, ok := os.LookupEnv("KAFKA_CLIENT_CERT_KEY")
	if !ok {
		return nil, errors.New("Kafka Client Certificate Key not found!")
	}

	clientCert, ok := os.LookupEnv("KAFKA_CLIENT_CERT")
	if !ok {
		return nil, errors.New("Kafka Client Certificate not found!")
	}

	roots := x509.NewCertPool()
	ok = roots.AppendCertsFromPEM([]byte(trustedCert))
	if !ok {
		return nil, errors.New("Unable to parse Root Cert. Please check your Heroku environment.")
	}
	// Setup certs for Sarama
	cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientCertKey))
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            roots,
	}, nil
}

// Extract the host:port pairs from the Kafka URL(s)
func brokerAddresses() ([]string, error) {
	URL := os.Getenv("KAFKA_URL")
	urls := strings.Split(URL, ",")
	addrs := make([]string, len(urls))
	for i, v := range urls {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		addrs[i] = u.Host
	}
	return addrs, nil
}
