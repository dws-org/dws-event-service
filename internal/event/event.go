package event

import (
	"errors"
	"time"
)

// Event ist ein minimales Beispiel-Modell.
type Event struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"` // RFC3339 erwartet
	Data      any    `json:"data,omitempty"`
}

// Validate prüft grundlegende Voraussetzungen des Events.
func (e *Event)Validate( ) error {
	if e == nil {
		return errors.New("event is nil")
	}
	if e.ID == "" {
		return errors.New("id is required")
	}
	if e.Type == "" {
		return errors.New("type is required")
	}
	if e.Timestamp == "" {
		// optional: akzeptieren, falls Timestamp optional sein soll; hier forciert
		return errors.New("timestamp is required")
	}
	if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
		return errors.New("timestamp must be RFC3339")
	}
	return nil
}