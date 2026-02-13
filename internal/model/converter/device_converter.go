package converter

import (
	"mertani_test/internal/entity"
	"mertani_test/internal/model"
)

func DeviceToResponse(device *entity.Device) *model.DeviceResponse {
	sensors := make([]model.SensorResponse, 0, len(device.Sensors))

	for _, sensor := range device.Sensors {
		sensors = append(sensors, *SensorToResponse(&sensor))
	}

	return &model.DeviceResponse{
		ID:        device.ID.String(),
		Name:      device.Name,
		Location:  device.Location,
		Status:    device.Status,
		Sensors:   sensors,
		CreatedAt: device.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: device.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
