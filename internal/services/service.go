package services

import services "github.com/mephistolie/chefbook-server/internal/repositories"

type Authorization interface {

}

type Recipes interface {

}

type Service struct {
	Authorization
	Recipes
}

func NewService(repos *services.Repository) *Service {
	return &Service{}
}