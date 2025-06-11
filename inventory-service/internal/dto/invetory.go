package dto

import "encoding/json"

type InventoryEvent struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"` // keep it generic for now
}
