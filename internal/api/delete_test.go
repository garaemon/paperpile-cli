package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTrashItem_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			return
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)

		resp := SyncResponse{SyncStartTime: 1234567890.0, SyncSession: "session-1"}
		w.Header().Set("Content-Type", "application/json")

		changes, ok := reqBody["clientChanges"].([]any)
		if !ok || len(changes) == 0 {
			json.NewEncoder(w).Encode(resp)
			return
		}

		change := changes[0].(map[string]any)
		if change["mcollection"] != "Library" {
			t.Errorf("mcollection = %v, want %q", change["mcollection"], "Library")
		}
		if change["action"] != "update" {
			t.Errorf("action = %v, want %q", change["action"], "update")
		}
		if change["id"] != "item-123" {
			t.Errorf("id = %v, want %q", change["id"], "item-123")
		}

		data, ok := change["data"].(map[string]any)
		if !ok {
			t.Fatal("expected data to be a map")
		}
		if data["trashed"] != float64(1) {
			t.Errorf("trashed = %v, want 1", data["trashed"])
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	err := client.TrashItem("item-123")
	if err != nil {
		t.Fatalf("TrashItem() error: %v", err)
	}
}

func TestTrashItem_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	client := newTestClient(server)
	err := client.TrashItem("item-123")
	if err == nil {
		t.Fatal("TrashItem() expected error for 500 response")
	}
}
