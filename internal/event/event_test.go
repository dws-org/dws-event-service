package event

import (
	"testing"
	"time"
)

func TestValidate_ValidEvent(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	ev := &Event{
		ID:        "evt-1",
		Type:      "user.created",
		Timestamp: now,
		Data:      map[string]any{"name": "Alice"},
	}

	if err := ev.Validate(); err != nil {
		t.Fatalf("expected valid event, got error: %v", err)
	}
}

func TestValidate_MissingID(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	ev := &Event{
		ID:        "",
		Type:      "user.created",
		Timestamp: now,
	}

	err := ev.Validate()
	if err == nil {
		t.Fatal("expected error for missing id, got nil")
	}
	if err.Error() != "id is required" {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestValidate_InvalidTimestamp(t *testing.T) {
	ev := &Event{
		ID:        "evt-2",
		Type:      "user.created",
		Timestamp: "not-a-timestamp",
	}

	err := ev.Validate()
	if err == nil {
		t.Fatal("expected error for invalid timestamp, got nil")
	}
	if err.Error() != "timestamp must be RFC3339" {
		t.Fatalf("unexpected error message: %v", err)
	}
}