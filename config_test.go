package heroku

import (
	"testing"
)

func TestEnvName(t *testing.T) {
	cases := []struct {
		TestName string
		Name     string
		Attr     string
		Result   string
	}{
		{"KAFKA_URL", "", "URL", "KAFKA_URL"},
		{"Named URL", "WHITE", "URL", "HEROKU_KAFKA_WHITE_URL"},

		// Guard against lower casing the name
		{"Lower cased name", "white", "url", "HEROKU_KAFKA_WHITE_URL"},
	}

	for _, test := range cases {
		t.Run(test.TestName, func(t *testing.T) {
			name := envName(test.Name, test.Attr)
			if name != test.Result {
				t.Errorf("envName returned %q but expected %q", name, test.Result)
			}
		})
	}
}
