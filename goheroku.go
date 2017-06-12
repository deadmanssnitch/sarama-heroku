package goheroku

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/joeshaw/envdecode"
	"log"
	"net/url"
	"strings"
)

// This library makes connecting to Heroku easy.
//
//

type AppConfig struct {
	URL           string `env:"KAFKA_URL,required"`
	TrustedCert   string `env:"KAFKA_TRUSTED_CERT,required"`
	ClientCertKey string `env:"KAFKA_CLIENT_CERT_KEY,required"`
	ClientCert    string `env:"KAFKA_CLIENT_CERT,required"`
	Prefix        string `env:"KAFKA_PREFIX"`
}

// <<Some helpful documentation here>>
//
func NewConsumer(topic string, consumerGroup string) (*cluster.Consumer, error) {
	ac, _ := setupConnection
	consumer, err := ac.createKafkaConsumer(brokerAddrs, tlsConfig)
	if err != nil {
		return nil, err
	}
}

func NewAsyncProducer() (*sarama.AsyncProducer, error) {
	ac, _ := setupConnection
	producer, err := ac.createKafkaAsyncProducer(brokerAddrs, tlsConfig)
	if err != nil {
		return nil, err
	}
}

func NewSyncProducer() (*sarama.SyncProducer, error) {
	ac, _ := setupConnection
	producer, err := ac.createKafkaSyncProducer(brokerAddrs, tlsConfig)
	if err != nil {
		return nil, err
	}
}

func setupConnection() (*AppConfig, error) {
	ac := AppConfig{}
	err := envdecode.Decode(&ac)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	tlsConfig, err := ac.createTLSConfig()
	if err != nil {
		return nil, err
	}
	brokerAddrs, err := ac.brokerAddresses()
	if err != nil {
		return nil, err
	}
}

// Create the TLS context, using the key and certificates provided.
func (ac *AppConfig) createTLSConfig() (*tls.Config, error) {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(ac.TrustedCert))
	if !ok {
		fmt.Printf("Unable to parse Root Cert:", ac.TrustedCert)
	}
	// Setup certs for Sarama
	cert, err := tls.X509KeyPair([]byte(ac.ClientCert), []byte(ac.ClientCertKey))
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            roots,
	}
	return tlsConfig, nil
}

// Extract the host:port pairs from the Kafka URL(s)
func (ac *AppConfig) brokerAddresses() ([]string, error) {
	urls := strings.Split(ac.URL, ",")
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

// Consumer group will default to Sarama if there is no value passed in
func (ac *AppConfig) createKafkaConsumer(brokers []string, tc *tls.Config) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Group.PartitionStrategy = cluster.StrategyRoundRobin
	config.ClientID = ac.ConsumerGroup
	config.Consumer.Return.Errors = true
	group := ac.ConsumerGroup
	if ac.Prefix != "" {
		group = strings.Join([]string{ac.Prefix, group}, "")
	}
	consumer, err := cluster.NewConsumer(brokers, group, []string{consumerTopic}, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (ac *AppConfig) createKafkaAsyncProducer(brokers []string, tc *tls.Config) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Default is WaitForLocal
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func (ac *AppConfig) createKafkaSyncProducer(brokers []string, tc *tls.Config) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Default is WaitForLocal
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	return producer, nil
}
