package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// newTestClient creates a Client that points to the given test server.
func newTestClient(server *httptest.Server) *Client {
	return &Client{
		session:    "test-session",
		httpClient: server.Client(),
		baseURL:    server.URL,
	}
}

func TestNewClient(t *testing.T) {
	client := NewClient("my-session")
	if client.session != "my-session" {
		t.Errorf("session = %q, want %q", client.session, "my-session")
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", client.baseURL, defaultBaseURL)
	}
}

func TestFetchCurrentUser_success(t *testing.T) {
	expectedUser := UserInfo{
		ID:          "user123",
		GoogleName:  "Test User",
		GoogleEmail: "test@example.com",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/users/me" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		origin := r.Header.Get("Origin")
		if origin != "https://app.paperpile.com" {
			t.Errorf("Origin header = %q, want %q", origin, "https://app.paperpile.com")
		}

		cookie, err := r.Cookie("plack_session")
		if err != nil {
			t.Errorf("missing plack_session cookie: %v", err)
		} else if cookie.Value != "test-session" {
			t.Errorf("cookie value = %q, want %q", cookie.Value, "test-session")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedUser)
	}))
	defer server.Close()

	client := newTestClient(server)
	user, err := client.FetchCurrentUser()
	if err != nil {
		t.Fatalf("FetchCurrentUser() error: %v", err)
	}
	if user.ID != expectedUser.ID {
		t.Errorf("user.ID = %q, want %q", user.ID, expectedUser.ID)
	}
	if user.GoogleName != expectedUser.GoogleName {
		t.Errorf("user.GoogleName = %q, want %q", user.GoogleName, expectedUser.GoogleName)
	}
}

func TestFetchCurrentUser_unauthorized(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("unauthorized"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.FetchCurrentUser()
	if err == nil {
		t.Fatal("FetchCurrentUser() expected error for 401 response")
	}
}
