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

//  NewClusterConsumer creates a bsm sarama-cluster consumer
//  Provide the topic, consumer group and a cluster config
func NewClusterConsumer(topic string, consumerGroup string, cfg *cluster.Config) (*cluster.Consumer, error) {
	tc, err := createTLSConfig()
	if err != nil {
		return nil, err
	}
	cfg.Net.TLS.Config = tc
	cfg.Net.TLS.Enable = true

	topic = AppendPrefixTo(topic)
	group := AppendPrefixTo(consumerGroup)
	brokers, err := brokerAddresses()
	if err != nil {
		return nil, err
	}

	consumer, err := cluster.NewConsumer(brokers, group, []string{topic}, cfg)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// A quick and easy cluster consumer built with default configs
func NewQuickClusterConsumer(topic string, groupID string) (*cluster.Consumer, error) {
	cfg := cluster.NewConfig()
	consumer, err := NewClusterConsumer(topic, groupID, cfg)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// TODO: Investigate use of/need for topic
func NewConsumer(cfg *sarama.Config) (sarama.Consumer, error) {
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

// A quick and easy sarama consumer built with default configs
func NewQuickConsumer() (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	consumer, err := NewConsumer(cfg)
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

// A quick and easy async producer built with default configs
func NewQuickAsyncProducer() (sarama.AsyncProducer, error) {
	cfg := sarama.NewConfig()
	producer, err := NewAsyncProducer(cfg)
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

// A quick and easy sync producer built with default configs
func NewQuickSyncProducer() (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := NewSyncProducer(cfg)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// Endangered function. TODO: Investigate usefulness of client
func NewClient(cfg *sarama.Config) (sarama.Client, error) {
	addrs, err := brokerAddresses()
	if err != nil {
		return nil, err
	}
	client, err := sarama.NewClient(addrs, cfg)
	if err != nil {
		return nil, err
	}
	return client, nil
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
	URL, ok := os.LookupEnv("KAFKA_URL")
	if !ok {
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
