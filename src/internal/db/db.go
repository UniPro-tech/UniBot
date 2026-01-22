package db

import (
	"unibot/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	dsn := "host=localhost user=root password=secret dbname=unibot port=5432 sslmode=disable"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func SetupDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&model.Guild{},
		&model.Member{},
		&model.AuditLogSetting{},
		&model.BotSystemSetting{},
		&model.PinSetting{},
		&model.RSSSetting{},
		&model.ScheduleSetting{},
		&model.TTSConnection{},
		&model.TTSPersonalSetting{},
	)
	return err
}
