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
	ConsumerGroup string `env:"KAFKA_CONSUMER_GROUP,default=agent-data"`
}

type KafkaClient struct {
	consumer *cluster.Consumer
	producer sarama.AsyncProducer
}

var consumerTopic string
var producerTopic string

func createKafkaClient() {
	ac := AppConfig{}
	err := envdecode.Decode(&ac)
	if err != nil {
		log.Fatal(err)
	}
	initTopics(&ac)
}

func initTopics(ac *AppConfig) {
	if ac.Prefix != "" {
		consumerTopic = strings.Join([]string{ac.Prefix, consumerTopic}, "")
		producerTopic = strings.Join([]string{ac.Prefix, producerTopic}, "")
	}
}

// Setup the Kafka client for producing and consuming messages.
// Use the specified configuration environment variables.
func createKafkaClient(ac *AppConfig) (*KafkaClient, error) {
	tlsConfig, err := ac.createTLSConfig()
	if err != nil {
		return nil, err
	}
	brokerAddrs, err := ac.brokerAddresses()
	if err != nil {
		return nil, err
	}

	fmt.Printf("All broker server certificates are valid!")

	consumer, err := ac.createKafkaConsumer(brokerAddrs, tlsConfig)
	if err != nil {
		return nil, err
	}
	producer, err := ac.createKafkaProducer(brokerAddrs, tlsConfig)
	if err != nil {
		return nil, err
	}
	fmt.Println("Kafka client created successfully.")
	return &KafkaClient{
		consumer: consumer,
		producer: producer,
	}, nil
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
	fmt.Println("TLS context created")
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
	fmt.Println("Broker addresses parsed.")
	return addrs, nil
}

func (ac *AppConfig) createKafkaConsumer(brokers []string, tc *tls.Config) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Group.PartitionStrategy = cluster.StrategyRoundRobin
	config.ClientID = ac.ConsumerGroup
	config.Consumer.Return.Errors = true
	fmt.Printf("Consuming topic %s on brokers: %s", consumerTopic, brokers)
	group := ac.ConsumerGroup
	if ac.Prefix != "" {
		group = strings.Join([]string{ac.Prefix, group}, "")
	}
	consumer, err := cluster.NewConsumer(brokers, group, []string{consumerTopic}, config)
	if err != nil {
		fmt.Printf("group: %s", group)
		fmt.Printf("client ID: %s", config.ClientID)
		return nil, err
	}
	fmt.Println("Kafka consumer created successfully")
	return consumer, nil
}

func (ac *AppConfig) createKafkaProducer(brokers []string, tc *tls.Config) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Producer.Return.Errors = true
	config.Producer.RequiredAcks = sarama.WaitForAll // Default is WaitForLocal
	config.ClientID = ac.ConsumerGroup
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	fmt.Println("Kafka producer created successfully")
	return producer, nil
}
