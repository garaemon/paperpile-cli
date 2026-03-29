package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFormatFirstAuthor_withLastName(t *testing.T) {
	item := &LibraryItem{
		Author: []Author{{First: "John", Last: "Doe", Formatted: "J. Doe"}},
	}
	got := item.FormatFirstAuthor()
	if got != "Doe" {
		t.Errorf("FormatFirstAuthor() = %q, want %q", got, "Doe")
	}
}

func TestFormatFirstAuthor_withFormattedOnly(t *testing.T) {
	item := &LibraryItem{
		Author: []Author{{Formatted: "Organization Inc."}},
	}
	got := item.FormatFirstAuthor()
	if got != "Organization Inc." {
		t.Errorf("FormatFirstAuthor() = %q, want %q", got, "Organization Inc.")
	}
}

func TestFormatFirstAuthor_noAuthors(t *testing.T) {
	item := &LibraryItem{}
	got := item.FormatFirstAuthor()
	if got != "" {
		t.Errorf("FormatFirstAuthor() = %q, want empty string", got)
	}
}

func TestFormatFirstAuthor_multipleAuthors(t *testing.T) {
	item := &LibraryItem{
		Author: []Author{
			{First: "Alice", Last: "Smith"},
			{First: "Bob", Last: "Jones"},
		},
	}
	got := item.FormatFirstAuthor()
	if got != "Smith" {
		t.Errorf("FormatFirstAuthor() = %q, want %q", got, "Smith")
	}
}

func TestFetchLibrary_success(t *testing.T) {
	items := []LibraryItem{
		{ID: "item1", Title: "Test Paper", Year: "2024"},
		{ID: "item2", Title: "Another Paper", Year: "2023"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/library" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	}))
	defer server.Close()

	client := newTestClient(server)
	got, err := client.FetchLibrary()
	if err != nil {
		t.Fatalf("FetchLibrary() error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("FetchLibrary() returned %d items, want 2", len(got))
	}
	if got[0].Title != "Test Paper" {
		t.Errorf("got[0].Title = %q, want %q", got[0].Title, "Test Paper")
	}
}

func TestFetchLibrary_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.FetchLibrary()
	if err == nil {
		t.Fatal("FetchLibrary() expected error for 500 response")
	}
}
