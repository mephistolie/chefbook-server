package response_body

import "github.com/google/uuid"

type Id struct {
	Id      uuid.UUID `json:"id"`
	Message string    `json:"message,omitempty"`
}
