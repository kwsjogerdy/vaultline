package healthcheck

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func makeServer(code int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	}))
}

func TestCheck_Healthy(t *testing.T) {
	srv := makeServer(http.StatusOK)
	defer srv.Close()

	c := New(srv.URL, 5*time.Second)
	s := c.Check()

	if !s.Reachable {
		t.Fatal("expected reachable")
	}
	if s.Sealed {
		t.Fatal("expected not sealed")
	}
	if s.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", s.StatusCode)
	}
	if s.Error != "" {
		t.Fatalf("unexpected error: %s", s.Error)
	}
}

func TestCheck_Sealed(t *testing.T) {
	srv := makeServer(503)
	defer srv.Close()

	c := New(srv.URL, 5*time.Second)
	s := c.Check()

	if !s.Reachable {
		t.Fatal("expected reachable")
	}
	if !s.Sealed {
		t.Fatal("expected sealed")
	}
	if s.Error == "" {
		t.Fatal("expected error message for sealed vault")
	}
}

func TestCheck_Unreachable(t *testing.T) {
	c := New("http://127.0.0.1:19999", 200*time.Millisecond)
	s := c.Check()

	if s.Reachable {
		t.Fatal("expected unreachable")
	}
	if s.Error == "" {
		t.Fatal("expected error message")
	}
}

func TestStatus_String_Healthy(t *testing.T) {
	s := Status{Reachable: true, StatusCode: 200, Latency: 3 * time.Millisecond}
	out := s.String()
	if out == "" {
		t.Fatal("expected non-empty string")
	}
}

func TestStatus_String_Unreachable(t *testing.T) {
	s := Status{Reachable: false, Error: "connection refused"}
	out := s.String()
	if out == "" {
		t.Fatal("expected non-empty string")
	}
}
