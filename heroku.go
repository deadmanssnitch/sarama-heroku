package heroku

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

//  NewClusterConsumer creates a bsm sarama-cluster consumer
//  Provide the topic, consumer group and a cluster config
func NewClusterConsumer(groupID string, topics []string, cfg *cluster.Config) (*cluster.Consumer, error) {
	if cfg == nil {
		cfg = cluster.NewConfig()
	}

	// Configure TLS from environment
	tc, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Net.TLS.Config = tc
	cfg.Net.TLS.Enable = true

	// Consumer groups require the Kafka prefix
	groupID = AppendPrefixTo(groupID)

	// Ensure all topics have the Kafka prefix applied
	for idx, topic := range topics {
		topics[idx] = AppendPrefixTo(topic)
	}

	brokers, err := brokerAddresses()
	if err != nil {
		return nil, err
	}

	consumer, err := cluster.NewConsumer(brokers, groupID, topics, cfg)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// TODO: Investigate use of/need for topic
func NewConsumer(cfg *sarama.Config) (sarama.Consumer, error) {
	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	tc, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Net.TLS.Config = tc
	cfg.Net.TLS.Enable = true

	brokers, err := brokerAddresses()
	if err != nil {
		return nil, err
	}
	consumer, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

//  NewAsyncProducer creates a sarama Async Producer
//  Provide a sarama config
//  For more information about the difference between Async and Sync producers
//  see the sarama documentation
func NewAsyncProducer(cfg *sarama.Config) (sarama.AsyncProducer, error) {
	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	tc, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Net.TLS.Config = tc
	cfg.Net.TLS.Enable = true

	brokers, err := brokerAddresses()
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

//  NewSyncProducer creates a sarama Sync Producer
//  Provide a sarama config
//  For more information about the difference between Async and Sync producers
//  see the sarama documentation
func NewSyncProducer(cfg *sarama.Config) (sarama.SyncProducer, error) {
	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	tc, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Net.TLS.Config = tc
	cfg.Net.TLS.Enable = true

	brokers, err := brokerAddresses()
	if err != nil {
		return nil, err
	}
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

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
	trustedCert := os.Getenv("KAFKA_TRUSTED_CERT")
	if trustedCert == "" {
		return nil, errors.New("Kafka Trusted Certificate not found!")
	}

	clientCertKey := os.Getenv("KAFKA_CLIENT_CERT_KEY")
	if clientCertKey == "" {
		return nil, errors.New("Kafka Client Certificate Key not found!")
	}

	clientCert := os.Getenv("KAFKA_CLIENT_CERT")
	if clientCert == "" {
		return nil, errors.New("Kafka Client Certificate not found!")
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(trustedCert))
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
	if URL == "" {
		return nil, errors.New("Kafka URL not found!")
	}
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
