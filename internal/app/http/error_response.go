package http

import (
	"context"
	"net/http"
	"os"

	appconfig "github.com/stratflow-labs/stratflow/internal/app/config"
	httpmiddleware "github.com/stratflow-labs/stratflow/internal/httpserver/middleware"
)

// TraceIDFromContext extracts request id from HTTP middleware context.
func TraceIDFromContext(ctx context.Context) string {
	if id, ok := httpmiddleware.RequestIDFromContext(ctx); ok {
		return id
	}
	return ""
}

// FinalizeErrorMessage returns a safe error message for API response.
// For non-500 statuses it keeps the provided fallback.
// For 500 it exposes internal error text only in non-production environments.
func FinalizeErrorMessage(status int, fallback string, err error) string {
	if status != http.StatusInternalServerError || err == nil {
		return fallback
	}
	if isProductionEnv() {
		return fallback
	}
	return err.Error()
}

func isProductionEnv() bool {
	return appconfig.IsProductionAppEnv(os.Getenv("APP_ENV"))
}
