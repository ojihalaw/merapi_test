package migration

import (
	"mertani_test/internal/entity"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Run(db *gorm.DB, log *logrus.Logger) {
	err := db.AutoMigrate(
		&entity.Device{},
		&entity.Sensor{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	log.Info("Migration success ✅")
}
