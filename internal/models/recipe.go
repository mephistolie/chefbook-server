package models

import (
	"github.com/golang/protobuf/ptypes/timestamp"
)

type Recipe struct {
	Id      int    `json:"-"`
	Name    string `json:"name"`
	OwnerId string `json:"owner_id"`

	Servings int16 `json:"servings"`
	Time     int16 `json:"time"`
	Calories int16 `json:"calories"`

	Ingredients []Selectable `json:"ingredients"`
	Cooking     []Selectable `json:"cooking"`

	Preview           string              `json:"preview"`
	Visibility        bool                `json:"visibility"`
	Encrypted         bool                `json:"encrypted"`
	CreationTimestamp timestamp.Timestamp `json:"creation_timestamp"`
	UpdateTimestamp   timestamp.Timestamp `json:"update_timestamp"`
}
