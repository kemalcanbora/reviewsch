package main

import (
	"log"
	"reviewsch/internal/app"
	"reviewsch/utils"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	err := utils.DefaultSystemRequirements().Validate()
	if err != nil {
		log.Fatal(err)
	}
}
