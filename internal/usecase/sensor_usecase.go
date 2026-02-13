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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SensorUseCase struct {
	DB                 *gorm.DB
	Log                *logrus.Logger
	Validator          *utils.Validator
	SensorRepository *repository.SensorRepository
}

func NewSensorUseCase(db *gorm.DB, logger *logrus.Logger, validator *utils.Validator,
	sensorRepository *repository.SensorRepository) *SensorUseCase {
	return &SensorUseCase{
		DB:                 db,
		Log:                logger,
		Validator:          validator,
		SensorRepository: sensorRepository,
	}
}

func (c *SensorUseCase) Create(ctx context.Context, request *model.CreateSensorRequest) error {
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

	exists, err := c.SensorRepository.ExistsByName(c.DB.WithContext(ctx), request.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("%w: %s", utils.ErrConflict, "sensor name already exist")
	}

	deviceUUID, err := uuid.Parse(request.DeviceID)
	if err != nil {
		return fmt.Errorf("%w: invalid device_id", utils.ErrValidation)
	}

	sensor := &entity.Sensor{
		DeviceID: deviceUUID,
		Name:     request.Name,
		Type:     request.Type,
		Unit:     request.Unit,
		IsActive: request.IsActive != nil && *request.IsActive,
	}

	if err := c.SensorRepository.Create(c.DB.WithContext(ctx), sensor); err != nil {
		c.Log.Warnf("Failed create sensor to database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *SensorUseCase) FindAll(ctx context.Context, pagination *utils.PaginationRequest) ([]model.SensorResponse, *utils.PaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var sensors []entity.Sensor
	total, err := c.SensorRepository.FindAll(c.DB.WithContext(ctx), &sensors, pagination)
	if err != nil {
		c.Log.Warnf("Failed find all sensor from database : %+v", err)
		return nil, nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	responses := make([]model.SensorResponse, len(sensors))
	for i, sensor := range sensors {
		responses[i] = *converter.SensorToResponse(&sensor)
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

func (c *SensorUseCase) FindByID(ctx context.Context, sensorID string) (*model.SensorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	sensor := &entity.Sensor{}
	_, err := c.SensorRepository.FindByIdWithDevice(
		c.DB.WithContext(ctx),
		sensor,
		sensorID,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Sensor not found, id=%s", sensorID)
			return nil, utils.ErrNotFound
		}
		c.Log.Warnf("Failed find sensor from database : %+v", err)
		return nil, fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}
	
	return converter.SensorToResponse(sensor), nil
}

func (c *SensorUseCase) Update(ctx context.Context, sensorID string, request *model.UpdateSensorRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	sensor := &entity.Sensor{}
	_, err := c.SensorRepository.FindById(c.DB.WithContext(ctx), sensor, sensorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Sensor not found, id=%s", sensorID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find sensor from database : %+v", err)
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
		sensor.Name = *request.Name
	}
	if request.Type != nil {
		sensor.Type = *request.Type
	}
	if request.Unit != nil {
		sensor.Unit = *request.Unit
	}
	if request.IsActive != nil {
		sensor.IsActive = *request.IsActive
	}

	err = c.SensorRepository.Update(c.DB.WithContext(ctx), sensor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Sensor not found, id=%s", sensorID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed update sensor from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}

func (c *SensorUseCase) Delete(ctx context.Context, sensorID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	
	sensor := &entity.Sensor{}
	_, err := c.SensorRepository.FindById(c.DB.WithContext(ctx), sensor, sensorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Sensor not found, id=%s", sensorID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed find sensor from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	err = c.SensorRepository.Delete(c.DB.WithContext(ctx), sensor)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Log.Infof("Sensor not found, id=%s", sensorID)
			return utils.ErrNotFound
		}
		c.Log.Warnf("Failed delete sensor from database : %+v", err)
		return fmt.Errorf("%w: %s", utils.ErrInternal, err.Error())
	}

	return nil
}
