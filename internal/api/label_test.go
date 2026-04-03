package api

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateLabel_success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sync" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			return
		}

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)

		resp := SyncResponse{SyncStartTime: 1234567890.0}
		w.Header().Set("Content-Type", "application/json")

		changes, ok := reqBody["clientChanges"].([]any)
		if !ok || len(changes) == 0 {
			json.NewEncoder(w).Encode(resp)
			return
		}

		change := changes[0].(map[string]any)

		if change["mcollection"] != "Collections" {
			t.Errorf("mcollection = %v, want %q", change["mcollection"], "Collections")
		}
		if change["action"] != "insert" {
			t.Errorf("action = %v, want %q", change["action"], "insert")
		}

		data := change["data"].(map[string]any)
		if data["cName"] != "NewTag" {
			t.Errorf("cName = %v, want %q", data["cName"], "NewTag")
		}
		if data["collectionType"] != "label" {
			t.Errorf("collectionType = %v, want %q", data["collectionType"], "label")
		}
		if data["cParent"] != "ROOT" {
			t.Errorf("cParent = %v, want %q", data["cParent"], "ROOT")
		}

		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	id, err := client.CreateLabel("NewTag")
	if err != nil {
		t.Fatalf("CreateLabel() error: %v", err)
	}
	if id == "" {
		t.Error("CreateLabel() returned empty ID")
	}
}

func TestCreateLabel_serverError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.CreateLabel("NewTag")
	if err == nil {
		t.Fatal("CreateLabel() expected error for 500 response")
	}
}
