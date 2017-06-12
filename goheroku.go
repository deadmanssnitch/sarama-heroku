package goheroku

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	"github.com/joeshaw/envdecode"
	"log"
	"net/url"
	"strings"
)

type appConfig struct {
	URL           string `env:"KAFKA_URL,required"`
	TrustedCert   string `env:"KAFKA_TRUSTED_CERT,required"`
	ClientCertKey string `env:"KAFKA_CLIENT_CERT_KEY,required"`
	ClientCert    string `env:"KAFKA_CLIENT_CERT,required"`
	Prefix        string `env:"KAFKA_PREFIX"`
	TLSConfig     *tls.Config
	BrokerAddrs   []string
}

// This function requires that you specify the topic that the returned consumer
// will consume from and the group that it will belong to.
func NewConsumer(topic string, consumerGroup string) (*cluster.Consumer, error) {
	ac, _ := setupConnection()
	consumer, err := ac.createKafkaConsumer(topic, consumerGroup, ac.BrokerAddrs, ac.TLSConfig)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// For move information about the difference between Async and Sync producers,
// see the Sarama docs
func NewAsyncProducer() (sarama.AsyncProducer, error) {
	ac, _ := setupConnection()
	producer, err := ac.createKafkaAsyncProducer(ac.BrokerAddrs, ac.TLSConfig)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// For move information about the difference between Async and Sync producers,
// see the Sarama docs
func NewSyncProducer() (sarama.SyncProducer, error) {
	ac, _ := setupConnection()
	producer, err := ac.createKafkaSyncProducer(ac.BrokerAddrs, ac.TLSConfig)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// To specify the topic or consumer group, the Kafka prefix needs
// to appended to it. This function makes it possible without having
// access to the app config.
func AppendPrefixTo(str string) string {
	ac := appConfig{}
	err := envdecode.Decode(&ac)
	if err != nil {
		log.Fatal(err)
	}
	if ac.Prefix != "" {
		str = strings.Join([]string{ac.Prefix, str}, "")
	}
	return str
}

func setupConnection() (*appConfig, error) {
	ac := appConfig{}
	err := envdecode.Decode(&ac)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	err = ac.createTLSConfig()
	if err != nil {
		return nil, err
	}
	err = ac.brokerAddresses()
	if err != nil {
		return nil, err
	}
	return &ac, nil
}

// Create the TLS context, using the key and certificates provided.
func (ac *appConfig) createTLSConfig() error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(ac.TrustedCert))
	if !ok {
		return errors.New("Unable to parse Root Cert. Please check your Heroku environment.")
	}
	// Setup certs for Sarama
	cert, err := tls.X509KeyPair([]byte(ac.ClientCert), []byte(ac.ClientCertKey))
	if err != nil {
		return err
	}
	ac.TLSConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
		RootCAs:            roots,
	}
	return nil
}

// Extract the host:port pairs from the Kafka URL(s)
func (ac *appConfig) brokerAddresses() error {
	urls := strings.Split(ac.URL, ",")
	addrs := make([]string, len(urls))
	for i, v := range urls {
		u, err := url.Parse(v)
		if err != nil {
			return err
		}
		addrs[i] = u.Host
		ac.BrokerAddrs = addrs
	}
	return nil
}

// Consumer group will default to Sarama if there is no value passed in
func (ac *appConfig) createKafkaConsumer(topic string, consumerGroup string, brokers []string, tc *tls.Config) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Net.TLS.Config = tc
	config.Net.TLS.Enable = true
	config.Group.PartitionStrategy = cluster.StrategyRoundRobin
	config.ClientID = consumerGroup
	config.Consumer.Return.Errors = true
	group := consumerGroup
	group = AppendPrefixTo(group)
	topic = AppendPrefixTo(topic)
	consumer, err := cluster.NewConsumer(brokers, group, []string{topic}, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func (ac *appConfig) createKafkaAsyncProducer(brokers []string, tc *tls.Config) (sarama.AsyncProducer, error) {
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

func (ac *appConfig) createKafkaSyncProducer(brokers []string, tc *tls.Config) (sarama.SyncProducer, error) {
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
