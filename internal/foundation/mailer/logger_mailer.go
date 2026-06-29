package mailer

import (
	"context"

	"github.com/stratflow-labs/stratflow/internal/foundation/logger"
)

// LoggerMailer logs outgoing messages instead of sending real emails.
type LoggerMailer struct{}

func (LoggerMailer) Send(_ context.Context, to, subject, body string) error {
	logger.Info("mailer: to=%s subject=%s body=%s", to, subject, body)

	return nil
}
