package http

import (
	"errors"
	"mertani_test/internal/model"
	"mertani_test/internal/usecase"
	"mertani_test/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type SensorController struct {
	Log     *logrus.Logger
	UseCase *usecase.SensorUseCase
}

func NewSensorController(useCase *usecase.SensorUseCase, logger *logrus.Logger) *SensorController {
	return &SensorController{
		Log:     logger,
		UseCase: useCase,
	}
}

// CreateSensor godoc
// @Summary Create Sensor
// @Description Create new sensor
// @Tags Sensors
// @Accept json
// @Produce json
// @Param request body model.CreateSensorRequest true "Sensor Request"
// @Success 200 {object} model.SensorResponse
// @Failure 400 {object} map[string]interface{}
// @Router /sensors [post]
func (c *SensorController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateSensorRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create sensor : %+v", err)
		switch {
		case errors.Is(err, utils.ErrValidation):
			return ctx.Status(fiber.StatusBadRequest).
				JSON(utils.ErrorResponse(fiber.StatusBadRequest, err.Error()))

		case errors.Is(err, utils.ErrConflict):
			return ctx.Status(fiber.StatusConflict).
				JSON(utils.ErrorResponse(fiber.StatusConflict, err.Error()))

		default: // internal error
			return ctx.Status(fiber.StatusInternalServerError).
				JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
		}
	}

	return ctx.Status(fiber.StatusCreated).
		JSON(utils.DefaultSuccessResponse(fiber.StatusCreated, "sensor created successfully"))
}

// FindAll godoc
// @Summary Get Sensors List
// @Description Get list of sensors with pagination
// @Tags Sensors
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param order_by query string false "Order by field"
// @Param sort_by query string false "Sort by direction (asc/desc)"
// @Param search query string false "Search term"
// @Success 200 {object} model.SensorResponse
// @Failure 500 {object} map[string]interface{}
// @Router /sensors [get]
func (c *SensorController) FindAll(ctx *fiber.Ctx) error {
	req := &utils.PaginationRequest{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", 10),
		OrderBy: ctx.Query("order_by", "created_at"),
		SortBy:  ctx.Query("sort_by", "desc"),
		Search:  ctx.Query("search", ""),
	}

	sensors, pagination, err := c.UseCase.FindAll(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponseWithPagination(fiber.StatusOK, "get list sensor successfully", sensors, pagination))
}

// FindByID godoc
// @Summary Get Sensor by ID
// @Description Get sensor details by ID
// @Tags Sensors
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Success 200 {object} model.SensorResponse
// @Failure 404 {object} map[string]interface{}
// @Router /sensors/{id} [get]
func (c *SensorController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	sensor, err := c.UseCase.FindByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "sensor not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "get detail sensor successfully", sensor))
}

// UpdateSensor godoc
// @Summary Update Sensor
// @Description Update sensor by ID
// @Tags Sensors
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Param request body model.UpdateSensorRequest true "Sensor Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /sensors/{id} [put]
func (c *SensorController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateSensorRequest)
	id := ctx.Params("id")

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.Update(ctx.Context(), id, request)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "sensor not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update sensor successfully"))
}

// DeleteSensor godoc
// @Summary Delete Sensor
// @Description Delete sensor by ID
// @Tags Sensors
// @Accept json
// @Produce json
// @Param id path string true "Sensor ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /sensors/{id} [delete]
func (c *SensorController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.UseCase.Delete(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "sensor not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "delete sensor successfully"))
}
