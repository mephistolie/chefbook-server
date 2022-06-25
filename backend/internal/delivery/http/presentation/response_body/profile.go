package response_body

import (
	"github.com/mephistolie/chefbook-server/internal/entity"
	"time"
)

type Username struct {
	Username string `json:"username"`
}

type MinimalProfileInfo struct {
	Id                int     `json:"id"`
	Username          *string `json:"username,omitempty"`
	CreationTimestamp string  `json:"creation_timestamp"`
	Avatar            *string `json:"avatar,omitempty"`
	IsPremium         bool    `json:"premium,omitempty"`
	Broccoins         int     `json:"broccoins"`
}

func NewMinimalProfileInfo(profile entity.ProfileInfo) MinimalProfileInfo {
	return MinimalProfileInfo{
		Id:                profile.Id,
		Username:          profile.Username,
		CreationTimestamp: profile.CreationTimestamp.UTC().Format(time.RFC3339),
		Avatar:            profile.Avatar,
		IsPremium:         profile.PremiumEndDate != nil && profile.PremiumEndDate.Unix() > time.Now().Unix(),
		Broccoins:         profile.Broccoins,
	}
}

func NewMinimalProfileInfoByProfile(profile entity.Profile) MinimalProfileInfo {
	return MinimalProfileInfo{
		Id:                profile.Id,
		Username:          profile.Username,
		CreationTimestamp: profile.CreationTimestamp.UTC().Format(time.RFC3339),
		Avatar:            profile.Avatar,
		IsPremium:         profile.PremiumEndDate != nil && profile.PremiumEndDate.Unix() > time.Now().Unix(),
		Broccoins:         profile.Broccoins,
	}
}

func NewUsersList(profiles []entity.ProfileInfo) []MinimalProfileInfo {
	users := make([]MinimalProfileInfo, len(profiles))
	for i, profile := range profiles {
		users[i] = NewMinimalProfileInfo(profile)
	}
	return users
}

type DetailedProfileInfo struct {
	Id                int       `json:"id"`
	Email             string    `json:"email"`
	Username          *string   `json:"username,omitempty"`
	CreationTimestamp time.Time `json:"creation_timestamp"`
	Avatar            *string   `json:"avatar,omitempty"`
	IsPremium         bool      `json:"premium,omitempty"`
	Broccoins         int       `json:"broccoins"`
	IsBlocked         bool      `json:"is_blocked,omitempty"`
}

func NewDetailedProfileInfo(profile entity.Profile) DetailedProfileInfo {
	return DetailedProfileInfo{
		Id:                profile.Id,
		Email:             profile.Email,
		Username:          profile.Username,
		CreationTimestamp: profile.CreationTimestamp.UTC(),
		Avatar:            profile.Avatar,
		IsPremium:         profile.PremiumEndDate != nil && profile.PremiumEndDate.Unix() > time.Now().Unix(),
		Broccoins:         profile.Broccoins,
		IsBlocked:         profile.IsBlocked,
	}
}
