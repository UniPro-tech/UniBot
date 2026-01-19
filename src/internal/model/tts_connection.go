package model

type TTSConnection struct {
	GuildID   string `gorm:"primaryKey;size:255"`
	ChannelID string `gorm:"not null;size:255;uniqueIndex:idx_guild_channel"`
	CreatedAt int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `gorm:"autoUpdateTime:nano"`
	Guild     Guild  `gorm:"foreignKey:GuildID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
