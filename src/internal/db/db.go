package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDB() (*pgxpool.Pool, error) {
	return pgxpool.New(
		context.Background(),
		"postgres://root:secret@localhost:5432/unibot",
	)
}
