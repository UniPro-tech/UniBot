package model

import "github.com/jackc/pgtype"

type BotSystemSettings struct {
	id    string       `gorm:"primaryKey;size:255"`
	value pgtype.JSONB `gorm:"not null"`
}
