package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GetItemLabels fetches the label IDs for a specific library item.
func (c *Client) GetItemLabels(itemID string) ([]string, error) {
	items, err := c.FetchLibrary()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch library: %w", err)
	}

	for _, item := range items {
		if item.ID == itemID {
			return item.LabelIDs, nil
		}
	}
	return nil, fmt.Errorf("item %s not found", itemID)
}

// AddLabel adds a label to a library item via the Sync API.
func (c *Client) AddLabel(itemID, labelID string) error {
	now := float64(time.Now().UnixMilli()) / 1000.0

	changes := []map[string]any{
		{
			"mcollection": "Library",
			"action":      "update",
			"id":          itemID,
			"timestamp":   now,
			"fields":      []string{"labels"},
			"push":        []string{labelID},
		},
	}

	_, err := c.pushSyncChanges(changes)
	if err != nil {
		return fmt.Errorf("failed to add label: %w", err)
	}
	return nil
}

// RemoveLabel removes a label from a library item via the Sync API.
func (c *Client) RemoveLabel(itemID, labelID string) error {
	now := float64(time.Now().UnixMilli()) / 1000.0

	changes := []map[string]any{
		{
			"mcollection": "Library",
			"action":      "update",
			"id":          itemID,
			"timestamp":   now,
			"fields":      []string{"labels"},
			"pull":        []string{labelID},
		},
	}

	_, err := c.pushSyncChanges(changes)
	if err != nil {
		return fmt.Errorf("failed to remove label: %w", err)
	}
	return nil
}

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

// DeleteLabel deletes a label by name via the Sync API.
func (c *Client) DeleteLabel(name string) error {
	labelID, err := c.ResolveLabelName(name)
	if err != nil {
		return err
	}

	return c.deleteCollection(labelID)
}

// ResolveLabelName finds a label ID by name. Returns the label ID or an error.
func (c *Client) ResolveLabelName(name string) (string, error) {
	labels, err := c.FetchLabels()
	if err != nil {
		return "", err
	}

	for _, label := range labels {
		if label.Name == name {
			return label.ID, nil
		}
	}
	return "", fmt.Errorf("label %q not found", name)
}

// GetItemLabelNames fetches the label names for a specific library item.
func (c *Client) GetItemLabelNames(itemID string) ([]string, error) {
	labelIDs, err := c.GetItemLabels(itemID)
	if err != nil {
		return nil, err
	}

	if len(labelIDs) == 0 {
		return nil, nil
	}

	labels, err := c.FetchLabels()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch labels: %w", err)
	}

	labelMap := make(map[string]string, len(labels))
	for _, label := range labels {
		labelMap[label.ID] = label.Name
	}

	var names []string
	for _, id := range labelIDs {
		if name, ok := labelMap[id]; ok {
			names = append(names, name)
		} else {
			names = append(names, id)
		}
	}
	return names, nil
}

// AssignLabelByName assigns a label to a library item, resolving the label by name.
func (c *Client) AssignLabelByName(itemID, labelName string) error {
	labelID, err := c.ResolveLabelName(labelName)
	if err != nil {
		return err
	}
	return c.AddLabel(itemID, labelID)
}

// UnassignLabelByName unassigns a label from a library item, resolving the label by name.
func (c *Client) UnassignLabelByName(itemID, labelName string) error {
	labelID, err := c.ResolveLabelName(labelName)
	if err != nil {
		return err
	}
	return c.RemoveLabel(itemID, labelID)
}
