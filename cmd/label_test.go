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
