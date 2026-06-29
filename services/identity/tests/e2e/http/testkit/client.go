package testkit

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

type Client struct {
	env   Env
	token string
}

type Response struct {
	StatusCode int
	Body       []byte
	JSON       map[string]any
}

func NewClient(env Env, token string) Client {
	return Client{env: env, token: token}
}

func (c Client) Get(t *testing.T, path string, query url.Values) Response {
	t.Helper()
	return c.DoJSON(t, http.MethodGet, path, query, nil)
}

func (c Client) Post(t *testing.T, path string, body any) Response {
	t.Helper()
	return c.DoJSON(t, http.MethodPost, path, nil, body)
}

func (c Client) Put(t *testing.T, path string, body any) Response {
	t.Helper()
	return c.DoJSON(t, http.MethodPut, path, nil, body)
}

func (c Client) Delete(t *testing.T, path string) Response {
	t.Helper()
	return c.DoJSON(t, http.MethodDelete, path, nil, nil)
}

func (c Client) DoJSON(t *testing.T, method, path string, query url.Values, body any) Response {
	t.Helper()

	var reader io.Reader = http.NoBody
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	reqURL := c.env.BaseURL + "/api" + path
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

	out := Response{StatusCode: resp.StatusCode, Body: payload}
	if len(bytes.TrimSpace(payload)) > 0 {
		if err := json.Unmarshal(payload, &out.JSON); err != nil {
			t.Fatalf("decode JSON response status=%d body=%s: %v", resp.StatusCode, payload, err)
		}
	}
	return out
}

func UserPath(id string) string {
	return fmt.Sprintf("/users/%s", url.PathEscape(id))
}
