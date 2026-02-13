package repository

import (
	"mertani_test/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceRepository struct {
	Repository[entity.Device]
	Log *logrus.Logger
}

func NewDeviceRepository(log *logrus.Logger) *DeviceRepository {
	return &DeviceRepository{
		Log: log,
	}
}

func (r *DeviceRepository) FindByIdWithSensors(db *gorm.DB, device *entity.Device, id any) (*entity.Device, error) {
	if err := db.Preload("Sensors").
		Where("id = ?", id).
		Take(device).Error; err != nil {
		return nil, err
	}
	return device, nil
}

func (r *DeviceRepository) CountByName(db *gorm.DB, name string) (int64, error) {
	var count int64
	err := db.Model(&entity.Device{}).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (r *DeviceRepository) ExistsByName(db *gorm.DB, name string) (bool, error) {
	count, err := r.CountByName(db, name)
	return count > 0, err
}
