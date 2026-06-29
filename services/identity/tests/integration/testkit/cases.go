package testkit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Case struct {
	Name       string
	Method     string
	Path       string
	Token      string
	WantStatus int
}

func RunCases(t *testing.T, router http.Handler, cases []Case) {
	t.Helper()

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), c.Method, c.Path, nil)
			if c.Token != "" {
				req.Header.Set("Authorization", c.Token)
			}
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)
			if rr.Code != c.WantStatus {
				t.Fatalf("status = %d, want %d", rr.Code, c.WantStatus)
			}
		})
	}
}
