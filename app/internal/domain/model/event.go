package model

type Event struct {
	ID          string                 `json:"id"`
	Payload     map[string]interface{} `json:"payload"`
	ProcessedAt string                 `json:"processed_at"`
}
