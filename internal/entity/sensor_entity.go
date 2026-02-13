package entity

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DeviceID  uuid.UUID `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"size:100;not null"`
	Type      string    `gorm:"size:50;not null"`
	Unit      string    `gorm:"size:20"`
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Device Device `gorm:"foreignKey:DeviceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}