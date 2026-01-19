package model

type Member struct {
	DiscordUserID string `gorm:"primaryKey;size:255"`
}
