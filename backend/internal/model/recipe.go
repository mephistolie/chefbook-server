package model

import (
	"time"
)

type Recipe struct {
	Id          int    `json:"id,omitempty" db:"recipe_id"`
	Name        string `json:"name"`
	OwnerId     int    `json:"owner_id,omitempty" db:"owner_id"`
	OwnerName   string `json:"owner_name,omitempty"`
	Owned       bool   `json:"owned,omitempty"`
	Language    string `json:"language,omitempty" db:"language"`
	Description string `json:"description,omitempty" db:"description"`
	Likes       int16  `json:"likes,omitempty" db:"likes"`

	Categories []Category `json:"categories,omitempty"`
	Favourite  bool       `json:"favourite,omitempty"`
	Liked      bool       `json:"liked,omitempty"`

	Servings       int16              `json:"servings,omitempty"`
	Time           int16              `json:"time,omitempty"`
	Calories       int16              `json:"calories,omitempty"`
	Macronutrients MacronutrientsInfo `json:"macronutrients,omitempty"`

	Ingredients interface{} `json:"ingredients"`
	Cooking     interface{} `json:"cooking"`

	Preview    string `json:"preview,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	Encrypted  bool   `json:"encrypted,omitempty"`

	CreationTimestamp time.Time `json:"creation_timestamp,omitempty" db:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp,omitempty" db:"update_timestamp"`
}

type RecipeInfo struct {
	Id        int    `json:"id,omitempty" db:"recipe_id"`
	Name      string `json:"name"`
	OwnerId   int    `json:"owner_id,omitempty" db:"owner_id"`
	OwnerName string `json:"owner_name,omitempty"`
	Owned     bool   `json:"owned,omitempty"`
	Language  string `json:"language,omitempty" db:"language"`
	Likes     int16  `json:"likes,omitempty" db:"likes"`

	Categories []Category `json:"categories,omitempty"`
	Favourite  bool       `json:"favourite,omitempty"`
	Liked      bool       `json:"liked,omitempty"`

	Servings int16 `json:"servings,omitempty"`
	Time     int16 `json:"time,omitempty"`

	Calories int16 `json:"calories,omitempty"`

	Preview    string `json:"preview,omitempty"`
	Visibility string `json:"visibility,omitempty"`
	Encrypted  bool   `json:"encrypted,omitempty"`

	CreationTimestamp time.Time `json:"creation_timestamp,omitempty" db:"creation_timestamp"`
	UpdateTimestamp   time.Time `json:"update_timestamp,omitempty" db:"update_timestamp"`
}

type RecipeCategoriesInput struct {
	RecipeId   int   `json:"recipe_id,omitempty" db:"recipe_id"`
	UserId     int   `json:"user_id,omitempty" db:"user_id"`
	Categories []int `json:"categories,omitempty"`
}

type FavouriteRecipeInput struct {
	RecipeId  int  `json:"recipe_id,omitempty" db:"recipe_id"`
	UserId    int  `json:"user_id,omitempty" db:"user_id"`
	Favourite bool `json:"favourite,omitempty"`
}

type RecipeLikeInput struct {
	RecipeId int  `json:"recipe_id,omitempty" db:"recipe_id"`
	UserId   int  `json:"user_id,omitempty" db:"user_id"`
	Liked    bool `json:"liked,omitempty"`
}

type RecipesRequestParams struct {
	UserId      int      `json:"user_id,omitempty"`
	AuthorId    int      `json:"author_id,omitempty"`
	Owned       bool     `json:"owned,omitempty"`
	Search      string   `json:"search,omitempty"`
	Page        int      `json:"last_recipe_id,omitempty"`
	PageSize    int      `json:"page_size,omitempty"`
	SortBy      string   `json:"sort_by,omitempty"`
	MinTime     int      `json:"min_time,omitempty"`
	MaxTime     int      `json:"max_time,omitempty"`
	MinCalories int      `json:"min_calories,omitempty"`
	MaxCalories int      `json:"max_calories,omitempty"`
	MinServings int      `json:"min_servings,omitempty"`
	MaxServings int      `json:"max_servings,omitempty"`
	Languages   []string `json:"language,omitempty"`
}
