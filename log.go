package htmltable

import (
	"context"
	"fmt"
	"log"
	"strings"
)

// Logger is a very simplistic structured logger, than should
// be overriden by integrations.
var Logger func(_ context.Context, msg string, fields ...any)

func init() {
	Logger = defaultLogger
}

var defaultLogger = func(_ context.Context, msg string, fields ...any) {
	var sb strings.Builder
	sb.WriteString(msg)
	if len(fields)%2 != 0 {
		panic(fmt.Errorf("number of logged fields is not even"))
	}
	for i := 0; i < len(fields); i += 2 {
		sb.WriteRune(' ')
		sb.WriteString(fmt.Sprint(fields[i]))
		sb.WriteRune('=')
		sb.WriteString(fmt.Sprint(fields[i+1]))
	}
	log.Print(sb.String())
}
