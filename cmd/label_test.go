package cmd

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/garaemon/paperpile/internal/api"
)

type mockLabelFetcher struct {
	labels []api.Collection
	err    error
}

func (m *mockLabelFetcher) FetchLabels() ([]api.Collection, error) {
	return m.labels, m.err
}

func TestExecLabelList_success(t *testing.T) {
	fetcher := &mockLabelFetcher{
		labels: []api.Collection{
			{ID: "id-1", Name: "ML", Count: 5},
			{ID: "id-2", Name: "Robotics", Count: 3},
		},
	}

	var buf bytes.Buffer
	err := execLabelList(fetcher, &buf)
	if err != nil {
		t.Fatalf("execLabelList() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ML") {
		t.Errorf("output should contain 'ML', got: %s", output)
	}
	if !strings.Contains(output, "Robotics") {
		t.Errorf("output should contain 'Robotics', got: %s", output)
	}
}

func TestExecLabelList_empty(t *testing.T) {
	fetcher := &mockLabelFetcher{labels: nil}

	var buf bytes.Buffer
	err := execLabelList(fetcher, &buf)
	if err != nil {
		t.Fatalf("execLabelList() error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "(no labels)") {
		t.Errorf("output should contain '(no labels)', got: %s", output)
	}
}

func TestExecLabelList_error(t *testing.T) {
	fetcher := &mockLabelFetcher{err: errors.New("api error")}

	var buf bytes.Buffer
	err := execLabelList(fetcher, &buf)
	if err == nil {
		t.Fatal("execLabelList() expected error")
	}
}

type mockLabelCreator struct {
	calledName string
	returnedID string
	err        error
}

func (m *mockLabelCreator) CreateLabel(name string) (string, error) {
	m.calledName = name
	return m.returnedID, m.err
}

func TestExecLabelCreate_success(t *testing.T) {
	creator := &mockLabelCreator{returnedID: "new-id-123"}

	var buf bytes.Buffer
	err := execLabelCreate(creator, &buf, "NewLabel")
	if err != nil {
		t.Fatalf("execLabelCreate() error: %v", err)
	}

	if creator.calledName != "NewLabel" {
		t.Errorf("calledName = %q, want %q", creator.calledName, "NewLabel")
	}

	output := buf.String()
	if !strings.Contains(output, "NewLabel") {
		t.Errorf("output should mention label name, got: %s", output)
	}
	if !strings.Contains(output, "new-id-123") {
		t.Errorf("output should mention label ID, got: %s", output)
	}
}

func TestExecLabelCreate_error(t *testing.T) {
	creator := &mockLabelCreator{err: errors.New("sync failed")}

	var buf bytes.Buffer
	err := execLabelCreate(creator, &buf, "NewLabel")
	if err == nil {
		t.Fatal("execLabelCreate() expected error")
	}
}
