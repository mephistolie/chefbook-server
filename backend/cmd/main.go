package main

import (
	"flag"

	"github.com/mephistolie/chefbook-server/internal/app"
)

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
	var configDir string
	flag.StringVar(&configDir, "configDir", "configs", "Path to directory with config files.")
	flag.Parse()

	app.Run(configDir)
}
