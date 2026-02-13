package main

import (
	"fmt"
	"mertani_test/internal/config"
	"mertani_test/internal/migration"
	"mertani_test/internal/utils"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Merapi IoT API
// @version 1.0
// @description API for managing devices and sensors
// @host localhost:8080
// @BasePath /api/v1
func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := utils.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	migration.Run(db, log)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	config.Bootstrap(&config.BootstrapConfig{
		DB:          db,
		App:         app,
		Log:         log,
		Validator:   validator,
		Config:      viperConfig,
	})

	webPort := viperConfig.GetInt("APP_PORT")
	err := app.Listen(fmt.Sprintf(":%d", webPort))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
