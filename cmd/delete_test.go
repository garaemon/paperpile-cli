package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type mockItemTrasher struct {
	trashedID string
	err       error
}

func (m *mockItemTrasher) TrashItem(itemID string) error {
	m.trashedID = itemID
	return m.err
}

func TestExecDelete_success(t *testing.T) {
	trasher := &mockItemTrasher{}

	var buf bytes.Buffer
	err := execDelete(trasher, &buf, "item-abc")
	if err != nil {
		t.Fatalf("execDelete() error: %v", err)
	}

	if trasher.trashedID != "item-abc" {
		t.Errorf("trashedID = %q, want %q", trasher.trashedID, "item-abc")
	}

	output := buf.String()
	if !strings.Contains(output, "item-abc") {
		t.Error("output should mention the item ID")
	}
	if !strings.Contains(output, "Done") {
		t.Error("output should confirm completion")
	}
}

func TestExecDelete_trashError(t *testing.T) {
	trasher := &mockItemTrasher{
		err: errors.New("sync failed"),
	}

	var buf bytes.Buffer
	err := execDelete(trasher, &buf, "item-abc")
	if err == nil {
		t.Fatal("execDelete() expected error")
	}
	if !strings.Contains(err.Error(), "delete failed") {
		t.Errorf("error = %q, want to contain 'delete failed'", err.Error())
	}
}
