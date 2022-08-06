package main

import (
	"flag"

	"chefbook-server/internal/app"
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
	var (
		configDir string
		useSMTP   bool
	)
	flag.StringVar(&configDir, "configDir", "configs", "Path to directory with config files.")
	flag.BoolVar(&useSMTP, "useSMTP", false, "If false, service doesn't use SMTP.")
	flag.Parse()

	app.Run(configDir, useSMTP)
}
