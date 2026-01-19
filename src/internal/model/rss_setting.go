package model

type RSSSetting struct {
	ID        string `gorm:"primaryKey;size:255"`
	URL       string `gorm:"not null"`
	ChannelID string `gorm:"not null"`
	Title     string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `gorm:"autoUpdateTime:nano"`
	GuildID   string `gorm:"not null"`
	Guild     Guild  `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
