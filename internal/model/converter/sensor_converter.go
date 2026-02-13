package converter

import (
	"mertani_test/internal/entity"
	"mertani_test/internal/model"
)

func SensorToResponse(sensor *entity.Sensor) *model.SensorResponse {
	return &model.SensorResponse{
		ID:        sensor.ID.String(),
		DeviceID:  sensor.DeviceID.String(),
		DeviceName: sensor.Device.Name,
		Name:      sensor.Name,
		Type:      sensor.Type,
		Unit:      sensor.Unit,
		IsActive:  sensor.IsActive,
		CreatedAt: sensor.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: sensor.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
