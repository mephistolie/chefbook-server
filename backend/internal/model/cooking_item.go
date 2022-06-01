package model

type CookingItem struct {
	Text     string   `json:"text"`
	Link     string   `json:"link,omitempty"`
	Time     int16    `json:"time,omitempty"`
	Pictures []string `json:"pictures,omitempty"`
	Type     string   `json:"type,omitempty"`
}
