package common_body

import "chefbook-server/internal/entity"

type Macronutrients struct {
	Protein       *int16 `json:"protein,omitempty"`
	Fats          *int16 `json:"fats,omitempty"`
	Carbohydrates *int16 `json:"carbohydrates,omitempty"`
}

func (m *Macronutrients) Entity() entity.Macronutrients {
	return entity.Macronutrients{
		Protein:       m.Protein,
		Fats:          m.Fats,
		Carbohydrates: m.Carbohydrates,
	}
}

func NewMacronutrients(macronutrients entity.Macronutrients) *Macronutrients {
	var macronutrientsPointer *Macronutrients = nil
	if macronutrients.Protein != nil || macronutrients.Fats != nil || macronutrients.Carbohydrates != nil {
		macronutrientsResponse := Macronutrients{
			Protein:       macronutrients.Protein,
			Fats:          macronutrients.Fats,
			Carbohydrates: macronutrients.Carbohydrates,
		}
		macronutrientsPointer = &macronutrientsResponse
	}

	return macronutrientsPointer
}
