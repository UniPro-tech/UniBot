package model

import "github.com/jackc/pgtype"

type BotSystemSetting struct {
	ID    string       `gorm:"primaryKey;size:255"`
	Value pgtype.JSONB `gorm:"not null"`
}
