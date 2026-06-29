package testkit

import (
	"fmt"
	"net/http"
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
	if got.JSON == nil {
		t.Fatalf("expected JSON error body for status %d", wantStatus)
	}
	if _, ok := got.JSON["code"].(string); !ok {
		nested, nestedOK := got.JSON["error"].(map[string]any)
		if _, codeOK := nested["code"].(string); !nestedOK || !codeOK {
			t.Fatalf("expected stable error code in body, got %v", got.JSON)
		}
	}
}

func RequireDataMap(t *testing.T, got Response) map[string]any {
	t.Helper()
	data, ok := got.JSON["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected response data object, got %v", got.JSON)
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

func RequireEqualString(t *testing.T, got, want, field string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s = %q, want %q", field, got, want)
	}
}

func RequireJSONContainsData(t *testing.T, got Response, field string) string {
	t.Helper()
	return RequireStringField(t, RequireDataMap(t, got), field)
}

func StatusName(status int) string {
	return fmt.Sprintf("%d %s", status, http.StatusText(status))
}
