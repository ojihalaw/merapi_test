package usecase

import (
	"context"
	"errors"
	"fmt"
	"mertani_test/internal/entity"
	"mertani_test/internal/model"
	"mertani_test/internal/model/converter"
	"mertani_test/internal/repository"
	"mertani_test/internal/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validator          *utils.Validator
	DeviceRepository *repository.DeviceRepository
}

func NewDeviceUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	deviceRepository *repository.DeviceRepository) *DeviceUseCase {
	return &DeviceUseCase{
		DB:                 db,
		Log:                logger,
		Validator:          validator,
		DeviceRepository: deviceRepository,
	}
}

func (c *DeviceUseCase) Create(ctx context.Context, request *model.CreateDeviceRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	exists, err := c.DeviceRepository.ExistsByName(c.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "device name already exist")
	}

	category := &entity.Device{
		Name: request.Name,
		Location: request.Location,
		Status: request.Status,
	}

	if err := c.DeviceRepository.Create(c.DB.WithContext(ctx), category); err != nil {
		c.Log.Warnf("Failed create device to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *DeviceUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.DeviceResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var devices []entity.Device
	total, err := c.DeviceRepository.FindAll(c.DB.WithContext(ctx), &devices, pagination)
	if err != nil {
		c.Log.Warnf("Failed find all device from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.DeviceResponse, len(devices))
	for i, device := range devices {
		responses[i] = *converter.DeviceToResponse(&device)
	}

	totalPage := int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))

	paginationRes := &utils.PaginationResponse{
		Page:      pagination.Page,
		Limit:     pagination.Limit,
		OrderBy:   pagination.OrderBy,
		SortBy:    pagination.SortBy,
		Search:    pagination.Search,
		TotalData: total,
		TotalPage: totalPage,
	}

	return responses, paginationRes, nil
}

func (c *DeviceUseCase) FindByID(ctx context.Context, deviceID string) (*model.DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	device := &entity.Device{} 

	_, err := c.DeviceRepository.FindByIdWithSensors(
		c.DB.WithContext(ctx),
		device,
		deviceID,
	)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Device not found, id=%s", deviceID)
			return nil, utils.ErrNotFound
		}
		c.Log.Warnf("Failed find device from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return converter.DeviceToResponse(device), nil
}

func (c *DeviceUseCase) Update(ctx context.Context, deviceID string, request *model.UpdateDeviceRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	device := &entity.Device{}
	_, err := c.DeviceRepository.FindById(c.DB.WithContext(ctx), device, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Device not found, id=%s", deviceID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find device from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.Validator.Validate.Struct(request)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {

			var messages []string
			for _, e := range validationErrors {
				messages = append(messages, e.Translate(c.Validator.Translator))
			}
			return fmt.Errorf("%w: %s", utils.ErrValidation, strings.Join(messages, ", "))
		}
		return fmt.Errorf("%w: %s", utils.ErrValidation, err.Error())
	}

	if request.Name != nil {
		device.Name = *request.Name
	}

	if request.Location != nil {
		device.Location = *request.Location
	}

	if request.Status != nil {
		device.Status = *request.Status
	}

	err = c.DeviceRepository.Update(c.DB.WithContext(ctx), device)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Device not found, id=%s", deviceID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find device from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *DeviceUseCase) Delete(ctx context.Context, deviceID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	device := &entity.Device{}
	_, err := c.DeviceRepository.FindById(c.DB.WithContext(ctx), device, deviceID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Device not found, id=%s", deviceID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find device from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.DeviceRepository.Delete(c.DB.WithContext(ctx), device)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Device not found, id=%s", deviceID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find device from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
