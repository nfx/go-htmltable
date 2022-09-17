package htmltable

import (
	"context"
	"testing"
)

func TestLogger(t *testing.T) {
	Logger(context.Background(), "message", "foo", "bar", "x", 1)
}

func TestLoggerNoFields(t *testing.T) {
	Logger(context.Background(), "message")
}

func TestLoggerWrongFields(t *testing.T) {
	defer func(){
		p := recover()
		if p == nil {
			t.Fatalf("there must be panic")
		}
	}()
	Logger(context.Background(), "message", 1)
}