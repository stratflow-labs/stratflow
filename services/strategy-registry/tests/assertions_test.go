package tests

import (
	"net/http"
	"slices"
	"testing"
)

func requireStatus(t *testing.T, got e2eHTTPResponse, want int) {
	t.Helper()
	if got.StatusCode != want {
		t.Fatalf("status = %d, want %d, body = %s", got.StatusCode, want, got.Body)
	}
}

func requireErrorStatus(t *testing.T, got e2eHTTPResponse, want int) {
	t.Helper()
	requireStatus(t, got, want)
	requireErrorBody(t, got, want)
}

func requireErrorStatusOneOf(t *testing.T, got e2eHTTPResponse, wants ...int) {
	t.Helper()
	if !slices.Contains(wants, got.StatusCode) {
		t.Fatalf("status = %d, want one of %v, body = %s", got.StatusCode, wants, got.Body)
	}
	requireErrorBody(t, got, got.StatusCode)
}

func requireErrorBody(t *testing.T, got e2eHTTPResponse, status int) {
	t.Helper()
	if got.JSON == nil {
		t.Fatalf("expected JSON error body for status %d", status)
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

func requireData(t *testing.T, got e2eHTTPResponse) map[string]any {
	t.Helper()
	data, ok := got.JSON["data"].(map[string]any)
	if !ok {
		t.Fatalf("expected data object, got %v", got.JSON)
	}
	return data
}

func requireDataString(t *testing.T, got e2eHTTPResponse, field string) string {
	t.Helper()
	data := requireData(t, got)
	value, ok := data[field].(string)
	if !ok || value == "" {
		t.Fatalf("expected non-empty data.%s in %v", field, data)
	}
	return value
}

func requireListTotalAtLeast(t *testing.T, got e2eHTTPResponse, minTotal float64) {
	t.Helper()
	requireStatus(t, got, http.StatusOK)
	data := requireData(t, got)
	total, ok := data["total"].(float64)
	if !ok {
		t.Fatalf("expected numeric total in %v", data)
	}
	if total < minTotal {
		t.Fatalf("total = %v, want >= %v", total, minTotal)
	}
}

func requireListTotal(t *testing.T, got e2eHTTPResponse, want float64) {
	t.Helper()
	requireStatus(t, got, http.StatusOK)
	data := requireData(t, got)
	total, ok := data["total"].(float64)
	if !ok {
		t.Fatalf("expected numeric total in %v", data)
	}
	if total != want {
		t.Fatalf("total = %v, want %v", total, want)
	}
}
