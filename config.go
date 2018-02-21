package heroku

import (
	"crypto/tls"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	prefix    string
	brokers   []string
	tlsConfig *tls.Config
}

// NewConfig creates a config based on settings in KAFKA_URL
func NewConfig() (*Config, error) {
	return NewConfigWithName("")
}

// NewConfigWithName returns a Config pulling from HEROKU_KAFKA_[NAME]
// environment variables. Using an empty string for name will use the unnamed
// Kafka instance.
func NewConfigWithName(name string) (*Config, error) {
	var issues = newError()
	var certsRequired = false

	cfg := &Config{
		brokers: make([]string, 0),
		prefix:  os.Getenv(envName(name, "PREFIX")),
	}

	// Validate the list of brokers
	brokers := strings.Split(os.Getenv(envName(name, "URL")), ",")
	for _, b := range brokers {
		uri, err := url.Parse(b)
		if err != nil {
			issues.add(envName(name, "URL"), err.Error())
			continue
		}

		if uri.Scheme != "kafka" && uri.Scheme != "kafka+ssl" {
			issues.add(envName(name, "URL"), "%q is an invalid scheme", uri.Scheme)
			continue
		}

		// Enable certificate validations if any of them require TLS
		if uri.Scheme == "kafka+ssl" {
			certsRequired = true
		}

		cfg.brokers = append(cfg.brokers, uri.Host)
	}

	// Validate certificates if any of the brokers specified kafka+ssl
	if certsRequired {
		clientCert := os.Getenv(envName(name, "CLIENT_CERT"))
		if clientCert == "" {
			issues.add(envName(name, "CLIENT_CERT"), "is required for tls")
		}

		clientKey := os.Getenv(envName(name, "CLIENT_CERT_KEY"))
		if clientKey == "" {
			issues.add(envName(name, "CLIENT_CERT_KEY"), "is required for tls")
		}

		trustedCert := os.Getenv(envName(name, "TRUSTED_CERT"))
		if trustedCert == "" {
			issues.add(envName(name, "TRUSTED_CERT"), "is required for tls")
		}

		tlsCfg, err := NewTLSConfig(trustedCert, clientCert, clientKey)
		if err != nil {
			issues.add("TLS certificates", err.Error())
		}

		cfg.tlsConfig = tlsCfg
	}

	// Create a helpful error message to aid in debugging
	if issues.any() {
		return nil, issues
	}

	return cfg, nil
}

// Brokers returns the list of Kafka brokers to connect to.
func (c *Config) Brokers() []string {
	return c.brokers
}

// TLS will be true if TLS is required for the connection.
func (c *Config) TLS() bool {
	return c.tlsConfig != nil
}

// TLSConfig returns the *tls.Config that is configured with the needed
// certificates for Kafka and custom certificate verification to work with
// Heroku's certificates.
func (c *Config) TLSConfig() *tls.Config {
	return c.tlsConfig
}

// Prefix is used to add the Heroku prefix to topics and consumer group ids. It
// is safe to use when no prefix is set or when the given name is already
// prefixed.
func (c *Config) Prefix(name string) string {
	prefix := c.prefix

	if strings.HasPrefix(name, prefix) {
		return name
	}

	return prefix + name
}

// envName returns the environment variable name for the given name and
// attribute combination. The default set (e.g. KAFKA_URL) is used when name is
// empty. When name is given the variables are in the format
// HEROKU_KAFKA_[NAME]_[ATTR].
func envName(name string, attr string) string {
	prefix := "KAFKA"

	if strings.TrimSpace(name) != "" {
		prefix = "HEROKU_KAFKA_" + name
	}

	return strings.ToUpper(prefix + "_" + attr)
}
