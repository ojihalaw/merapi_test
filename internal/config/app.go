package config

import (
	"mertani_test/internal/delivery/http"
	"mertani_test/internal/delivery/http/route"
	"mertani_test/internal/repository"
	"mertani_test/internal/usecase"
	"mertani_test/internal/utils"

	_ "mertani_test/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Log         *logrus.Logger
	Validator   *utils.Validator
	Config      *viper.Viper
}

func Bootstrap(config *BootstrapConfig) {
	deviceRepository := repository.NewDeviceRepository(config.Log)
	deviceUseCase := usecase.NewDeviceUseCase(config.DB, config.Log, config.Validator, deviceRepository)
	deviceController := http.NewDeviceController(deviceUseCase, config.Log)	

	sensorRepository := repository.NewSensorRepository(config.Log)
	sensorUseCase := usecase.NewSensorUseCase(config.DB, config.Log, config.Validator, sensorRepository)
	sensorController := http.NewSensorController(sensorUseCase, config.Log)
	
	routeConfig := route.RouteConfig{
		App:                config.App,
		DeviceController: deviceController,
		SensorController: sensorController,
	}
	routeConfig.Setup()
}
