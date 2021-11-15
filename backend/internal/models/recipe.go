package models

import (
	"time"
)

type Recipe struct {
	Id      int    `json:"id,omitempty" db:"recipe_id"`
	Name    string `json:"name"`
	OwnerId int    `json:"owner_id,omitempty" db:"owner_id"`

	Favourite bool `json:"favourite,omitempty"`
	Liked     bool `json:"liked,omitempty"`

	Servings int16 `json:"servings,omitempty"`
	Time     int16 `json:"time,omitempty"`
	Calories int16 `json:"calories,omitempty"`

	Ingredients interface{} `json:"ingredients"`
	Cooking     interface{} `json:"cooking"`

	Preview           string    `json:"preview,omitempty"`
	Visibility        string    `json:"visibility,omitempty"`
	Encrypted         bool      `json:"encrypted,omitempty"`
	CreationTimestamp time.Time `json:"creation_timestamp,omitempty" db:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp,omitempty" db:"update_timestamp"`
}

type FavouriteRecipeInput struct {
	RecipeId        int  `json:"recipe_id,omitempty" db:"recipe_id"`
	Favourite bool `json:"favourite,omitempty"`
}
