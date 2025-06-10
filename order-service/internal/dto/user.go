package dto

import "encoding/json"

type UserEvent struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"` // keep it generic for now
}
