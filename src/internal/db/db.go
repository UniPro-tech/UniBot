package db

import (
	"os"
	"unibot/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	// dsn := "postgres://root:secret@localhost:5432/unibot?sslmode=disable"
	dsn := os.Getenv("PG_DSN")
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
		&model.TTSDictionary{},
	)
	db.Migrator().DropColumn(&model.TTSPersonalSetting{}, "speaker_seed")
	db.Migrator().DropColumn(&model.TTSPersonalSetting{}, "speaker_pitch")
	db.Migrator().DropColumn(&model.TTSPersonalSetting{}, "speed_scale")
	return err
}
