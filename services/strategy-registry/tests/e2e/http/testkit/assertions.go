package testkit

import (
	"net/http"
	"slices"
	"testing"
)

func RequireStatus(t *testing.T, got Response, want int) {
	t.Helper()
	if got.StatusCode != want {
		t.Fatalf("status = %d, want %d, body = %s", got.StatusCode, want, got.Body)
	}
}

func RequireErrorCode(t *testing.T, got Response, wantStatus int) {
	t.Helper()
	RequireStatus(t, got, wantStatus)
	RequireErrorBody(t, got)
}

func RequireErrorCodeOneOf(t *testing.T, got Response, wantStatuses ...int) {
	t.Helper()
	if !slices.Contains(wantStatuses, got.StatusCode) {
		t.Fatalf("status = %d, want one of %v, body = %s", got.StatusCode, wantStatuses, got.Body)
	}
	RequireErrorBody(t, got)
}

func RequireErrorBody(t *testing.T, got Response) {
	t.Helper()
	if got.JSON == nil {
		t.Fatalf("expected JSON error body, got status=%d body=%s", got.StatusCode, got.Body)
	}
	if _, ok := got.JSON["code"].(string); ok {
		return
	}
	if nested, ok := got.JSON["error"].(map[string]any); ok {
		if _, ok := nested["code"].(string); ok {
			return
		}
	}
	if _, ok := got.JSON["error_message"].(string); ok {
		return
	}
	t.Fatalf("expected stable error code/message in body, got %v", got.JSON)
}

func RequireDataMap(t *testing.T, got Response) map[string]any {
	t.Helper()
	data, ok := got.JSON["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %v", got.JSON)
	}
	return data
}

func RequireStringField(t *testing.T, data map[string]any, field string) string {
	t.Helper()
	value, ok := data[field].(string)
	if !ok || value == "" {
		t.Fatalf("expected non-empty string field %q in %v", field, data)
	}
	return value
}

func RequireDataString(t *testing.T, got Response, field string) string {
	t.Helper()
	return RequireStringField(t, RequireDataMap(t, got), field)
}

func RequireEqualString(t *testing.T, got, want, field string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %q, want %q", field, got, want)
	}
}

func RequireListTotal(t *testing.T, got Response, want float64) {
	t.Helper()
	RequireStatus(t, got, http.StatusOK)
	total, ok := RequireDataMap(t, got)["total"].(float64)
	if !ok {
		t.Fatalf("expected numeric total in %v", got.JSON)
	}
	if total != want {
		t.Fatalf("total = %v, want %v", total, want)
	}
}

func RequireListTotalAtLeast(t *testing.T, got Response, minTotal float64) {
	t.Helper()
	RequireStatus(t, got, http.StatusOK)
	total, ok := RequireDataMap(t, got)["total"].(float64)
	if !ok {
		t.Fatalf("expected numeric total in %v", got.JSON)
	}
	if total < minTotal {
		t.Fatalf("total = %v, want >= %v", total, minTotal)
	}
}
