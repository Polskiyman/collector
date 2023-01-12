package main

import (
	"log"

	"collector/internal/app"
	"collector/pkg"
)

const configPath = "config.json"

func main() {
	config, err := pkg.ParseFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	a := app.NewApp(config)
	a.Run()
}
