package heroku_test

import (
	"errors"
	"os"
	"reflect"
	"testing"

	heroku "github.com/deadmanssnitch/sarama-heroku"
)

func TestAppendPrefixTo(t *testing.T) {
	cases := map[string]struct {
		Prefix string
		Input  string
		Output string
	}{
		"adds KAFKA_PREFIX": {
			Prefix: "michigan-98173.",
			Input:  "events",
			Output: "michigan-98173.events",
		},
		"ignores empty prefix": {
			Prefix: "",
			Input:  "widgets",
			Output: "widgets",
		},
		"does not double append prefixes": {
			Prefix: "ohio-89287.",
			Input:  "ohio-89287.farms",
			Output: "ohio-89287.farms",
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := os.Setenv("KAFKA_PREFIX", test.Prefix)
			if err != nil {
				t.Fatal("unable to set KAFKA_PREFIX", err)
			}
			defer func() { _ = os.Unsetenv("KAFKA_PREFIX") }()

			result := heroku.AppendPrefixTo(test.Input)

			if result != test.Output {
				t.Errorf("%q did not match expected result of %q", result, test.Output)
			}
		})
	}
}

func TestBrokers(t *testing.T) {
	cases := map[string]struct {
		Env     string
		Err     error
		Brokers []string
	}{
		"simple": {
			Env:     "kafka+ssl://ec2-00-000-00-000.aws.com:9096,kafka+ssl://ec2-00-000-00-001.aws.com:9098",
			Err:     nil,
			Brokers: []string{"ec2-00-000-00-000.aws.com:9096", "ec2-00-000-00-001.aws.com:9098"},
		},
		"environment not set": {
			Env:     "",
			Err:     errors.New("KAFKA_URL is not set in environment"),
			Brokers: nil,
		},
		"invalid url": {
			Env:     "ec2-00-000-00-000.aws.com:9096",
			Err:     errors.New("kafka urls should start with kafka+ssl://"),
			Brokers: nil,
		},
	}

	for name, test := range cases {
		t.Run(name, func(t *testing.T) {
			err := os.Setenv("KAFKA_URL", test.Env)
			if err != nil {
				t.Fatal("unable to set KAFKA_URL", err)
			}
			defer func() { _ = os.Unsetenv("KAFKA_URL") }()

			result, err := heroku.Brokers()
			if test.Err == nil && err != nil {
				t.Error("expected no error but received", err)
			}

			if err != nil && err.Error() != test.Err.Error() {
				t.Errorf("error of %q did not match expected of %q", err, test.Err)
			}

			if !reflect.DeepEqual(result, test.Brokers) {
				t.Errorf("%q did not match expected result of %q", result, test.Brokers)
			}
		})
	}
}
