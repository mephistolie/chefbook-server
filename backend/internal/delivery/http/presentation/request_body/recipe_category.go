package request_body

import "github.com/google/uuid"

type RecipeCategoriesInput struct {
	Categories []uuid.UUID `json:"categories"`
}
