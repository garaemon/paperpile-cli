package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

type mockFileAttacher struct {
	calledItemID  string
	calledPath    string
	attachmentID  string
	err           error
}

func (m *mockFileAttacher) AttachFile(itemID, filePath string) (string, error) {
	m.calledItemID = itemID
	m.calledPath = filePath
	return m.attachmentID, m.err
}

func TestExecAttach_success(t *testing.T) {
	attacher := &mockFileAttacher{
		attachmentID: "att-456",
	}

	var buf bytes.Buffer
	err := execAttach(attacher, &buf, "pub-123", "/path/to/supplement.pdf")
	if err != nil {
		t.Fatalf("execAttach() error: %v", err)
	}

	if attacher.calledItemID != "pub-123" {
		t.Errorf("calledItemID = %q, want %q", attacher.calledItemID, "pub-123")
	}
	if attacher.calledPath != "/path/to/supplement.pdf" {
		t.Errorf("calledPath = %q, want %q", attacher.calledPath, "/path/to/supplement.pdf")
	}

	output := buf.String()
	if !strings.Contains(output, "supplement.pdf") {
		t.Error("output should contain file name")
	}
	if !strings.Contains(output, "pub-123") {
		t.Error("output should contain item ID")
	}
	if !strings.Contains(output, "att-456") {
		t.Error("output should contain attachment ID")
	}
}

func TestExecAttach_attachError(t *testing.T) {
	attacher := &mockFileAttacher{
		err: errors.New("permission denied"),
	}

	var buf bytes.Buffer
	err := execAttach(attacher, &buf, "pub-123", "/path/to/file.pdf")
	if err == nil {
		t.Fatal("execAttach() expected error")
	}
	if !strings.Contains(err.Error(), "attach failed") {
		t.Errorf("error = %q, want to contain 'attach failed'", err.Error())
	}
}
