package heroku

import (
	"crypto/tls"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

// NewClusterConsumer creates a github.com/bsm/sarama-cluster.Consumer based on
// Heroku Kafka standard environment configs. Giving nil for cfg will create a
// generic config.
func NewClusterConsumer(groupID string, topics []string, cfg *cluster.Config) (*cluster.Consumer, error) {
	herokuCfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = cluster.NewConfig()
	}

	cfg.Net.TLS.Enable = herokuCfg.TLS()
	cfg.Net.TLS.Config = herokuCfg.TLSConfig()

	// Consumer groups require the Kafka prefix
	groupID = herokuCfg.Prefix(groupID)

	// Ensure all topics have the Kafka prefix applied
	for idx, topic := range topics {
		topics[idx] = herokuCfg.Prefix(topic)
	}

	return cluster.NewConsumer(herokuCfg.Brokers(), groupID, topics, cfg)
}

// NewConsumer creates a github.com/Shopify/sarama.Consumer configured from the
// standard Heroku Kafka environment.
func NewConsumer(cfg *sarama.Config) (sarama.Consumer, error) {
	herokuCfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	cfg.Net.TLS.Enable = herokuCfg.TLS()
	cfg.Net.TLS.Config = herokuCfg.TLSConfig()

	consumer, err := sarama.NewConsumer(herokuCfg.Brokers(), cfg)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

// NewAsyncProducer creates a github.com/Shopify/sarama.AsyncProducer
// configured from the standard Heroku Kafka environment. When publishing
// messages to Multitenant Kafka all topics need to start with KAFKA_PREFIX
// which is best added using AppendPrefixTo.
func NewAsyncProducer(cfg *sarama.Config) (sarama.AsyncProducer, error) {
	herokuCfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	cfg.Net.TLS.Enable = herokuCfg.TLS()
	cfg.Net.TLS.Config = herokuCfg.TLSConfig()

	return sarama.NewAsyncProducer(herokuCfg.Brokers(), cfg)
}

// NewSyncProducer creates a github.com/Shopify/sarama.SyncProducer configured
// from the standard Heroku Kafka environment. When publishing messages to
// Multitenant Kafka all topics need to start with KAFKA_PREFIX which is best
// added using AppendPrefixTo.
func NewSyncProducer(cfg *sarama.Config) (sarama.SyncProducer, error) {
	herokuCfg, err := NewConfig()
	if err != nil {
		return nil, err
	}

	if cfg == nil {
		cfg = sarama.NewConfig()
	}

	cfg.Net.TLS.Enable = herokuCfg.TLS()
	cfg.Net.TLS.Config = herokuCfg.TLSConfig()

	return sarama.NewSyncProducer(herokuCfg.Brokers(), cfg)
}

// AppendPrefixTo adds the env variable KAFKA_PREFIX to the given string if
// necessary. Heroku requires prefixing topics and consumer group names with
// the prefix on multi-tenant plans. It is safe to use on dedicated clusters if
// KAFKA_PREFIX is not set.
func AppendPrefixTo(name string) string {
	prefix := os.Getenv("KAFKA_PREFIX")

	if strings.HasPrefix(name, prefix) {
		return name
	}

	return prefix + name
}

// Create the TLS context, using the key and certificates provided.
func TLSConfig() (*tls.Config, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	return config.TLSConfig(), nil
}

// Brokers returns a list of host:port addresses for the Kafka brokers set in
// KAFKA_URL.
func Brokers() ([]string, error) {
	config, err := NewConfig()
	if err != nil {
		return nil, err
	}

	return config.Brokers(), nil
}
