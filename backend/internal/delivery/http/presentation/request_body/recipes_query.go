package request_body

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/entity"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"strings"
)

type RecipesQuery struct {
	AuthorId    *uuid.UUID
	Owned       bool
	Saved       bool
	Search      *string
	Page        int
	PageSize    int
	SortBy      string
	Languages   *[]string
	MinTime     *int
	MaxTime     *int
	MinServings *int
	MaxServings *int
	MinCalories *int
	MaxCalories *int
}

func (p *RecipesQuery) Validate(userId uuid.UUID) error {
	if p.Owned && p.AuthorId != nil && *p.AuthorId != userId {
		return failure.InvalidBody
	}

	if p.Owned {
		p.AuthorId = &userId
	}

	if p.Search != nil && *p.Search == "" {
		p.Search = nil
	}

	if p.Page == 0 {
		p.Page = 1
	}

	if p.Page < 0 {
		return failure.InvalidBody
	}

	if p.PageSize == 0 {
		p.PageSize = 10
	}

	if p.PageSize < 0 {
		return failure.InvalidBody
	}

	if p.PageSize > 50 {
		p.PageSize = 50
	}

	if p.SortBy == "" {
		p.SortBy = entity.SortingCreationTimestamp
	}
	p.SortBy = strings.ToLower(p.SortBy)

	switch p.SortBy {
	case entity.SortingCreationTimestamp, entity.SortingUpdateTimestamp, entity.SortingLikes, entity.SortingTime,
		entity.SortingServings, entity.SortingCalories:
	default:
		return failure.InvalidBody
	}

	if p.MinTime != nil && *p.MinTime <= 0 {
		p.MinTime = nil
	}

	if p.MaxTime != nil && *p.MaxTime <= 0 {
		p.MaxTime = nil
	}

	if p.MinServings != nil && *p.MinServings <= 0 {
		p.MinServings = nil
	}

	if p.MaxServings != nil && *p.MaxServings <= 0 {
		p.MaxServings = nil
	}

	if p.MinCalories != nil && *p.MinCalories <= 0 {
		p.MinCalories = nil
	}

	if p.MaxCalories != nil && *p.MaxCalories <= 0 {
		p.MaxCalories = nil
	}

	if p.MinTime != nil && p.MaxTime != nil && *p.MinTime > *p.MaxTime {
		return failure.InvalidBody
	}

	if p.MinServings != nil && p.MaxServings != nil && *p.MinServings > *p.MaxServings {
		return failure.InvalidBody
	}

	if p.MinCalories != nil && p.MaxCalories != nil && *p.MinCalories > *p.MaxCalories {
		return failure.InvalidBody
	}

	return nil
}

func (p *RecipesQuery) Entity() entity.RecipesQuery {
	return entity.RecipesQuery{
		AuthorId:    p.AuthorId,
		Saved:       p.Saved,
		Search:      p.Search,
		Page:        p.Page,
		PageSize:    p.PageSize,
		SortBy:      p.SortBy,
		Languages:   p.Languages,
		MinTime:     p.MinTime,
		MaxTime:     p.MaxTime,
		MinCalories: p.MinCalories,
		MaxCalories: p.MaxCalories,
		MinServings: p.MinServings,
		MaxServings: p.MaxServings,
	}
}
