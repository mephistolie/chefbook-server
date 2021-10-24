package models

import (
	"time"
)

type Recipe struct {
	Id      int    `json:"-"`
	Name    string `json:"name"`
	OwnerId int    `json:"owner_id"`

	Servings int16 `json:"servings"`
	Time     int16 `json:"time"`
	Calories int16 `json:"calories"`

	Ingredients []Selectable `json:"ingredients"`
	Cooking     []Selectable `json:"cooking"`

	Preview           string    `json:"preview"`
	Visibility        string    `json:"visibility"`
	Encrypted         bool      `json:"encrypted"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp"`
}
