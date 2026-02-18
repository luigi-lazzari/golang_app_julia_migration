package service

import (
	"testing"
)

func TestMapExternalToInternal(t *testing.T) {
	external := []NewsExternalItem{
		{ID: "1", Description: "Desc 1", Channel: "A"},
		{ID: "2", Description: "Desc 2", Channel: "B"},
	}

	internal := MapExternalToInternal(external)

	if len(internal) != len(external) {
		t.Fatalf("Expected length %d, got %d", len(external), len(internal))
	}

	for i, item := range internal {
		if item.ID != external[i].ID {
			t.Errorf("Expected ID %s, got %s", external[i].ID, item.ID)
		}
		if item.Description != external[i].Description {
			t.Errorf("Expected Description %s, got %s", external[i].Description, item.Description)
		}
		if item.Channel != external[i].Channel {
			t.Errorf("Expected Channel %s, got %s", external[i].Channel, item.Channel)
		}
	}
}

type MockExternalGateway struct {
	data []NewsExternalItem
	err  error
}

func (m *MockExternalGateway) GetNotificationPreferences() ([]NewsExternalItem, error) {
	return m.data, m.err
}

type MockInternalGateway struct {
	lastData []NewsItem
	err      error
}

func (m *MockInternalGateway) UpdateNotificationNewsPreferences(news []NewsItem) error {
	m.lastData = news
	return m.err
}

func TestOrchestrateNotificationPreferencesUpdate(t *testing.T) {
	mockExt := &MockExternalGateway{
		data: []NewsExternalItem{{ID: "1", Description: "D", Channel: "C"}},
	}
	mockInt := &MockInternalGateway{}

	svc := NewOrchestratorService(mockExt, mockInt)

	err := svc.OrchestrateNotificationPreferencesUpdate()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(mockInt.lastData) != 1 {
		t.Fatalf("Expected 1 item updated, got %d", len(mockInt.lastData))
	}
}
