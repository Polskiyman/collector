package main

import (
	"collector/pkg"
	"log"

	"collector/internal/app"
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
