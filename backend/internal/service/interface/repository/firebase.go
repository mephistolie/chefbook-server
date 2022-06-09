package repository

import "github.com/mephistolie/chefbook-server/internal/entity"

type FirebaseMigration interface {
	GetProfile(credentials entity.Credentials) (entity.FirebaseProfile, error)
	GetProfileData(profile entity.FirebaseProfile) (entity.FirebaseUserData, error)
}
