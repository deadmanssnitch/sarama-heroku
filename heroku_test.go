package heroku_test

import (
	"os"
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

			result := heroku.AppendPrefixTo(test.Input)

			// Test clean up
			err = os.Unsetenv("KAFKA_PREFIX")
			if err != nil {
				t.Log("unable to clear KAFKA_PREFIX", err)
			}

			if result != test.Output {
				t.Errorf("%q did not match expected result of %q", result, test.Output)
			}
		})
	}
}
