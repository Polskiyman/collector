package main

import (
	"log"

	"collector/internal"
	"collector/internal/app"
)

const configPath = "config.json"

func main() {
	config, err := internal.ParseFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(config)
	a.Run()
}
