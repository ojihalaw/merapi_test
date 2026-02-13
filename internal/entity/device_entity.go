package entity

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `gorm:"size:100;not null"`
	Location  string    `gorm:"size:150"`
	Status    string    `gorm:"size:50;default:'active'"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Sensors []Sensor `gorm:"foreignKey:DeviceID"`
}