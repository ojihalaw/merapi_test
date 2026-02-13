package main

import (
	"fmt"
	"mertani_test/internal/config"
	"mertani_test/internal/migration"
	"mertani_test/internal/utils"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogger(viperConfig)
	db := config.NewDatabase(viperConfig, log)
	validator := utils.NewValidator(viperConfig)
	app := config.NewFiber(viperConfig)

	migration.Run(db, log)

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
