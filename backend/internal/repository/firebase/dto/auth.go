package dto

import "chefbook-server/internal/entity"

type FirebaseCredentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewFirebaseCredentials(credentials entity.Credentials) *FirebaseCredentials {
	return &FirebaseCredentials{
		Email:    credentials.Email,
		Password: credentials.Password,
	}
}

type FirebaseProfile struct {
	IdToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalId      string `json:"localId"`
	Registered   bool   `json:"registered"`
}

func (p *FirebaseProfile) Entity() entity.FirebaseProfile {
	return entity.FirebaseProfile{
		IdToken:      p.IdToken,
		Email:        p.Email,
		RefreshToken: p.RefreshToken,
		ExpiresIn:    p.ExpiresIn,
		LocalId:      p.LocalId,
		Registered:   p.Registered,
	}
}
