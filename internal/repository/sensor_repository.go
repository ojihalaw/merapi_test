package repository

import (
	"mertani_test/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SensorRepository struct {
	Repository[entity.Sensor]
	Log *logrus.Logger
}

func NewSensorRepository(log *logrus.Logger) *SensorRepository {
	return &SensorRepository{
		Log: log,
	}
}

func (r *SensorRepository) FindByIdWithDevice(db *gorm.DB, sensor *entity.Sensor, id any) (*entity.Sensor, error) {
	if err := db.Preload("Device").
		Where("id = ?", id).
		Take(sensor).Error; err != nil {
		return nil, err
	}
	return sensor, nil
}

func (r *SensorRepository) CountByName(db *gorm.DB, name string) (int64, error) {
	var count int64
	err := db.Model(&entity.Sensor{}).Where("name = ?", name).Count(&count).Error
	return count, err
}

func (r *SensorRepository) ExistsByName(db *gorm.DB, name string) (bool, error) {
	count, err := r.CountByName(db, name)
	return count > 0, err
}
