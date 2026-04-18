package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vaultline/vaultline/internal/notify"
)

func makeServer(t *testing.T, status int, fn func(body map[string]string)) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m map[string]string
		_ = json.NewDecoder(r.Body).Decode(&m)
		if fn != nil {
			fn(m)
		}
		w.WriteHeader(status)
	}))
}

func TestSend_Success(t *testing.T) {
	var got map[string]string
	srv := makeServer(t, 200, func(b map[string]string) { got = b })
	defer srv.Close()

	n := notify.New(notify.Config{WebhookURL: srv.URL, Channel: "#ops", Timeout: 2 * time.Second})
	if err := n.Send("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["text"] != "hello" {
		t.Errorf("expected text=hello, got %q", got["text"])
	}
	if got["channel"] != "#ops" {
		t.Errorf("expected channel=#ops, got %q", got["channel"])
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	srv := makeServer(t, 500, nil)
	defer srv.Close()

	n := notify.New(notify.Config{WebhookURL: srv.URL})
	if err := n.Send("msg"); err == nil {
		t.Fatal("expected error for 500 status")
	}
}

func TestSend_MissingURL(t *testing.T) {
	n := notify.New(notify.Config{})
	if err := n.Send("msg"); err == nil {
		t.Fatal("expected error when URL is empty")
	}
}

func TestSyncMessage_Success(t *testing.T) {
	msg := notify.SyncMessage("production", 3, 1, 0, nil)
	if msg == "" {
		t.Fatal("expected non-empty message")
	}
	for _, want := range []string{"production", "+3", "~1", "-0"} {
		if !contains(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
	}
}

func TestSyncMessage_Error(t *testing.T) {
	msg := notify.SyncMessage("staging", 0, 0, 0, fmt.Errorf("vault unreachable"))
	if !contains(msg, "FAILED") {
		t.Errorf("expected FAILED in message, got %q", msg)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
