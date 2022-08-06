package repository

import (
	firebaseRepo "chefbook-server/internal/repository/firebase"
	"chefbook-server/internal/repository/postgres"
	"chefbook-server/internal/repository/s3"
	"chefbook-server/internal/service/interface/repository"
	firebase "firebase.google.com/go/v4"
	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type Repository struct {
	Auth            repository.Auth
	Profile         repository.Profile
	Recipe          repository.Recipe
	RecipeOwnership repository.RecipeOwnership
	RecipeSharing   repository.RecipeSharing
	Encryption      repository.Encryption
	Category        repository.Category
	ShoppingList    repository.ShoppingList
	File            repository.File
	Migration       repository.FirebaseMigration
}

func NewRepository(db *sqlx.DB, client *minio.Client, firebaseApp *firebase.App, firebaseApiKey string) *Repository {
	var migrationRepo repository.FirebaseMigration = nil
	if firebaseApp != nil {
		migrationRepo = firebaseRepo.NewMigrationRepo(*firebaseApp, firebaseApiKey)
	}

	return &Repository{
		Auth:            postgres.NewAuthPostgres(db),
		Profile:         postgres.NewProfilePostgres(db),
		RecipeOwnership: postgres.NewRecipeOwnershipPostgres(db),
		Recipe:          postgres.NewRecipePostgres(db),
		RecipeSharing:   postgres.NewRecipeSharingPostgres(db),
		Encryption:      postgres.NewEncryptionPostgres(db),
		Category:        postgres.NewCategoryPostgres(db),
		ShoppingList:    postgres.NewShoppingListPostgres(db),
		File:            s3.NewAWSFileManager(client),
		Migration:       migrationRepo,
	}
}
