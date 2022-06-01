package model

type MacronutrientsInfo struct {
	Protein       int16 `json:"protein,omitempty" db:"protein"`
	Fats          int16 `json:"fats,omitempty" db:"fats"`
	Carbohydrates int16 `json:"carbohydrates" db:"carbohydrates"`
}
