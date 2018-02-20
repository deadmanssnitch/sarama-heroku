package heroku

import (
	"bytes"
	"fmt"
)

type Error struct {
	messages map[string][]string
}

func (e *Error) Error() string {
	out := &bytes.Buffer{}

	fmt.Fprintln(out, "Kafka config has the following errors:")
	for attr, issues := range e.messages {
		for _, message := range issues {
			fmt.Fprintf(out, "  - %s: %s\n", attr, message)
		}
	}

	return out.String()
}

func newError() *Error {
	return &Error{
		messages: make(map[string][]string),
	}
}

func (e *Error) any() bool {
	return len(e.messages) > 0
}

func (e *Error) add(attr string, format string, attrs ...interface{}) {
	messages, ok := e.messages[attr]
	if !ok {
		messages = make([]string, 0)
	}

	e.messages[attr] = append(messages, fmt.Sprintf(format, attrs...))
}
