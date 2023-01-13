package main

import (
	"log"

	"collector/internal/app"
	"collector/pkg/config"
)

const configPath = "config.json"

func main() {
	config, err := config.ParseFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(config)
	a.Run()
}
