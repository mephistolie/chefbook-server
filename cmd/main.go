package main

import "github.com/mephistolie/chefbook-server/internal/app"

const configDir = "configs"

func main() {
	app.Run(configDir)
}
