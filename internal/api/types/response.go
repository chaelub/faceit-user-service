package types

type (
	Response struct {
		Success bool        `json:"success"`
		Error   string      `json:"error,omitempty"`
		Payload interface{} `json:"payload"`
	}
)
