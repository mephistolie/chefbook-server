package main

import (
	"github.com/mephistolie/chefbook-server/internal/app"
)

const configDir = "configs"

// @title ChefBook API
// @version 1.0
// @description ChefBook API Server

// @contact.name   ChefBook API Support
// @contact.email  support@chefbook.space

// @host api.chefbook.space
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app.Run(configDir)
}
