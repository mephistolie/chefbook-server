package entity

import "github.com/google/uuid"

const (
	SortingCreationTimestamp = "creation_timestamp"
	SortingUpdateTimestamp   = "update_timestamp"
	SortingLikes             = "likes"
	SortingTime              = "time"
	SortingServings          = "servings"
	SortingCalories          = "calories"
)

type RecipesQuery struct {
	AuthorId    *uuid.UUID
	Saved       bool
	Search      *string
	Page        int
	PageSize    int
	SortBy      string
	MinTime     *int
	MaxTime     *int
	MinCalories *int
	MaxCalories *int
	MinServings *int
	MaxServings *int
	Languages   *[]string
}
