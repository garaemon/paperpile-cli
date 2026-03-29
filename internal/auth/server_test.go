package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleCallback_postSession(t *testing.T) {
	sessionCh := make(chan string, 1)

	req := httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader("my-session-value"))
	w := httptest.NewRecorder()

	handleCallback(w, req, sessionCh)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	select {
	case session := <-sessionCh:
		if session != "my-session-value" {
			t.Errorf("session = %q, want %q", session, "my-session-value")
		}
	default:
		t.Fatal("expected session to be sent to channel")
	}
}

func TestHandleCallback_optionsReturnsOK(t *testing.T) {
	sessionCh := make(chan string, 1)

	req := httptest.NewRequest(http.MethodOptions, "/callback", nil)
	w := httptest.NewRecorder()

	handleCallback(w, req, sessionCh)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	cors := resp.Header.Get("Access-Control-Allow-Origin")
	if cors != "*" {
		t.Errorf("CORS header = %q, want %q", cors, "*")
	}
}

func TestHandleCallback_getMethodNotAllowed(t *testing.T) {
	sessionCh := make(chan string, 1)

	req := httptest.NewRequest(http.MethodGet, "/callback", nil)
	w := httptest.NewRecorder()

	handleCallback(w, req, sessionCh)

	resp := w.Result()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusMethodNotAllowed)
	}
}

func TestHandleCallback_emptySession(t *testing.T) {
	sessionCh := make(chan string, 1)

	req := httptest.NewRequest(http.MethodPost, "/callback", strings.NewReader(""))
	w := httptest.NewRecorder()

	handleCallback(w, req, sessionCh)

	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusBadRequest)
	}
}

func TestHandleSetupPage_containsBookmarklet(t *testing.T) {
	w := httptest.NewRecorder()
	handleSetupPage(w, 18080)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		t.Errorf("Content-Type = %q, want text/html", contentType)
	}

	body := w.Body.String()
	if !strings.Contains(body, "localhost:18080") {
		t.Error("setup page should contain the callback URL with the correct port")
	}
	if !strings.Contains(body, "Paperpile CLI Login") {
		t.Error("setup page should contain the title")
	}
}
