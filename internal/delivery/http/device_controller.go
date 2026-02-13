package http

import (
	"errors"
	"mertani_test/internal/model"
	"mertani_test/internal/usecase"
	"mertani_test/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type DeviceController struct {
	Log     *logrus.Logger
	UseCase *usecase.DeviceUseCase
}

func NewDeviceController(useCase *usecase.DeviceUseCase, logger *logrus.Logger) *DeviceController {
	return &DeviceController{
		Log:     logger,
		UseCase: useCase,
	}
}

// CreateDevice godoc
// @Summary Create Device
// @Description Create new device
// @Tags Devices
// @Accept json
// @Produce json
// @Param request body model.CreateDeviceRequest true "Device Request"
// @Success 200 {object} model.DeviceResponse
// @Failure 400 {object} map[string]interface{}
// @Router /devices [post]
func (c *DeviceController) Create(ctx *fiber.Ctx) error {
	request := new(model.CreateDeviceRequest)

	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(utils.ErrorResponse(fiber.StatusBadRequest, "Failed to parse request body"))
	}

	err = c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.Warnf("Failed to create device : %+v", err)
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
		JSON(utils.DefaultSuccessResponse(fiber.StatusCreated, "device created successfully"))
}

// FindAll godoc
// @Summary Get List of Devices
// @Description Get all devices with pagination
// @Tags Devices
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param order_by query string false "Field to order by"
// @Param sort_by query string false "Sort direction (asc or desc)"
// @Param search query string false "Search term"
// @Success 200 {object} model.DeviceResponse
// @Failure 500 {object} map[string]interface{}
// @Router /devices [get]
func (c *DeviceController) FindAll(ctx *fiber.Ctx) error {
	req := &utils.PaginationRequest{
		Page:    ctx.QueryInt("page", 1),
		Limit:   ctx.QueryInt("limit", 10),
		OrderBy: ctx.Query("order_by", "created_at"),
		SortBy:  ctx.Query("sort_by", "desc"),
		Search:  ctx.Query("search", ""),
	}

	devices, pagination, err := c.UseCase.FindAll(ctx.Context(), req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, err.Error()))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponseWithPagination(fiber.StatusOK, "get list device successfully", devices, pagination))
}

// FindByID godoc
// @Summary Get Device by ID
// @Description Get device details by ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} model.DeviceResponse
// @Failure 404 {object} map[string]interface{}
// @Router /devices/{id} [get]
func (c *DeviceController) FindByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	device, err := c.UseCase.FindByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "device not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.SuccessResponse(fiber.StatusOK, "get detail device successfully", device))
}

// Update godoc
// @Summary Update Device
// @Description Update device details by ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Param request body model.UpdateDeviceRequest true "Update Device Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /devices/{id} [put]
func (c *DeviceController) Update(ctx *fiber.Ctx) error {
	request := new(model.UpdateDeviceRequest)
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
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "device not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "update device successfully"))
}

// Delete godoc
// @Summary Delete Device
// @Description Delete device by ID
// @Tags Devices
// @Accept json
// @Produce json
// @Param id path string true "Device ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /devices/{id} [delete]
func (c *DeviceController) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	err := c.UseCase.Delete(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).
				JSON(utils.ErrorResponse(fiber.StatusNotFound, "device not found"))
		}

		return ctx.Status(fiber.StatusInternalServerError).
			JSON(utils.ErrorResponse(fiber.StatusInternalServerError, "internal server error"))
	}

	return ctx.Status(fiber.StatusOK).
		JSON(utils.DefaultSuccessResponse(fiber.StatusOK, "delete device successfully"))
}
