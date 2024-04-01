package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/paimon_bank/internal/config"
)

func main() {
	viperConfig := config.NewViper()
	app := config.NewFiber()
	db := config.NewDatabase(viperConfig)
	log := config.NewLogger()
	aws := config.NewAws()
	validate := validator.New()
	config.Bootstrap(&config.BootstrapConfig{
		App:      app,
		DB:       db,
		Log:      log,
		Aws:      aws,
		Validate: validate,
	})

	err := app.Listen(fmt.Sprintf(":%d", 8080))
	if err != nil {
		log.Fatal("Failed to start server: %w \n", err)
	}
}
