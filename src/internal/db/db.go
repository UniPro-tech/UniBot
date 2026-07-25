package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	// dsn := "postgres://root:secret@localhost:5432/unibot?sslmode=disable"
	dsn := os.Getenv("PG_DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
