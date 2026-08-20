package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewDB は DB 接続を開く。
// ロガーは gorm.Open に渡すことで、接続時に gorm 自身が出すログも捕捉できる。
func NewDB(l gormlogger.Interface) (*gorm.DB, error) {
	// dsn := "postgres://root:secret@localhost:5432/unibot?sslmode=disable"
	dsn := os.Getenv("PG_DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: l})
}
