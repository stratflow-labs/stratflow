package migrations

import (
	"fmt"
	"os"
	"strings"
)

type serviceGooseLogger struct {
	label string
}

func newServiceGooseLogger(service, label string) serviceGooseLogger {
	_ = service
	return serviceGooseLogger{label: strings.TrimSpace(label)}
}

func (l serviceGooseLogger) Printf(format string, v ...any) {
	fmt.Println(l.formatGooseMessage(format, v...))
}

func (l serviceGooseLogger) Fatalf(format string, v ...any) {
	fmt.Fprintln(os.Stderr, l.formatGooseMessage(format, v...))
	os.Exit(1)
}

func (l serviceGooseLogger) formatGooseMessage(format string, v ...any) string {
	msg := strings.TrimSpace(fmt.Sprintf(format, v...))
	if strings.Contains(msg, "goose: no migrations to run. current version:") {
		if l.label == "seeds" {
			msg = strings.Replace(msg, "goose: no migrations to run. current version:", "No seeds to run. current version:", 1)
		} else {
			msg = strings.Replace(msg, "goose: no migrations to run. current version:", "No migrations to run. current version:", 1)
		}
	}
	return msg
}
