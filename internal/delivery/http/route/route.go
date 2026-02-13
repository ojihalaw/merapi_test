package route

import (
	"mertani_test/internal/delivery/http"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App                *fiber.App
	DeviceController *http.DeviceController
	SensorController *http.SensorController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	api := c.App.Group("/api/v1")

	device := api.Group("/devices")
	device.Post("", c.DeviceController.Create)
	device.Get("", c.DeviceController.FindAll)
	device.Get("/:id", c.DeviceController.FindByID)
	device.Put("/:id", c.DeviceController.Update)
	device.Delete("/:id", c.DeviceController.Delete)

	sensor := api.Group("/sensors")
	sensor.Post("", c.SensorController.Create)
	sensor.Get("", c.SensorController.FindAll)
	sensor.Get("/:id", c.SensorController.FindByID)
	sensor.Put("/:id", c.SensorController.Update)
	sensor.Delete("/:id", c.SensorController.Delete)
	
}