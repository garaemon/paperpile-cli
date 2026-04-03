package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateLabel creates a new label via the Sync API.
func (c *Client) CreateLabel(name string) (string, error) {
	now := float64(time.Now().UnixMilli()) / 1000.0
	labelID := uuid.New().String()

	changes := []map[string]any{
		{
			"mcollection": "Collections",
			"action":      "insert",
			"id":          labelID,
			"timestamp":   now,
			"data": map[string]any{
				"_id":            labelID,
				"cName":          name,
				"cParent":        "ROOT",
				"cSortOrder":     -1,
				"cStyle":         "0",
				"cHidden":        0,
				"collectionType": "label",
				"created":        now,
				"updated":        now,
			},
		},
	}

	_, err := c.pushSyncChanges(changes)
	if err != nil {
		return "", fmt.Errorf("failed to create label: %w", err)
	}
	return labelID, nil
}
