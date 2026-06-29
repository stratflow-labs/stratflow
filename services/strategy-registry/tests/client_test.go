package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type strategyE2EClient struct {
	env   strategyE2EEnv
	token string
}

type e2eHTTPResponse struct {
	StatusCode int
	Body       []byte
	JSON       map[string]any
}

func newStrategyE2EClient(env strategyE2EEnv, token string) strategyE2EClient {
	return strategyE2EClient{env: env, token: token}
}

func (c strategyE2EClient) get(t *testing.T, path string, query url.Values) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodGet, c.env.BaseURL+"/api"+path, query, nil)
}

func (c strategyE2EClient) post(t *testing.T, path string, body any) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodPost, c.env.BaseURL+"/api"+path, nil, body)
}

func (c strategyE2EClient) patch(t *testing.T, path string, body any) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodPatch, c.env.BaseURL+"/api"+path, nil, body)
}

func (c strategyE2EClient) delete(t *testing.T, path string) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodDelete, c.env.BaseURL+"/api"+path, nil, nil)
}

func (c strategyE2EClient) identityPost(t *testing.T, env strategyE2EEnv, path string, body any) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodPost, env.IdentityURL+"/api"+path, nil, body)
}

func (c strategyE2EClient) identityDelete(t *testing.T, env strategyE2EEnv, path string) e2eHTTPResponse {
	t.Helper()
	return c.doJSON(t, http.MethodDelete, env.IdentityURL+"/api"+path, nil, nil)
}

func (c strategyE2EClient) doJSON(t *testing.T, method, reqURL string, query url.Values, body any) e2eHTTPResponse {
	t.Helper()

	var reader io.Reader = http.NoBody
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}

	req, err := http.NewRequestWithContext(context.Background(), method, reqURL, reader)
	if err != nil {
		t.Fatalf("build request %s %s: %v", method, reqURL, err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if strings.TrimSpace(c.token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(c.token))
	}

	resp, err := c.env.HTTPClient.Do(req)
	if err != nil {
		t.Fatalf("send request %s %s: %v", method, reqURL, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("close response body: %v", err)
		}
	}()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}
	out := e2eHTTPResponse{StatusCode: resp.StatusCode, Body: payload}
	if len(bytes.TrimSpace(payload)) > 0 {
		if err := json.Unmarshal(payload, &out.JSON); err != nil {
			t.Fatalf("decode JSON response status=%d body=%s: %v", resp.StatusCode, payload, err)
		}
	}
	return out
}

func strategyPath(id string) string {
	return fmt.Sprintf("/strategies/%s", url.PathEscape(id))
}

func attributePath(strategyID, attributeID string) string {
	return fmt.Sprintf("%s/attributes/%s", strategyPath(strategyID), url.PathEscape(attributeID))
}

func valuePath(strategyID, attributeID, valueID string) string {
	return fmt.Sprintf("%s/values/%s", attributePath(strategyID, attributeID), url.PathEscape(valueID))
}
