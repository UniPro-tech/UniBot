package model

type ScheduleSetting struct {
	ID        string `gorm:"primaryKey;size:255"`
	ChannelID string `gorm:"not null;size:255;index:idx_channel_guild"`
	Content   string `gorm:"type:text;not null"`
	NextRunAt int64  `gorm:"not null"`
	Cron      string `gorm:"not null"`
	CreatedAt int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt int64  `gorm:"autoUpdateTime:nano"`
	GuildID   string `gorm:"not null;size:255;index:idx_channel_guild"`
	AuthorID  string `gorm:"not null;size:255"`
	Author    Member `gorm:"foreignKey:AuthorID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Guild     Guild  `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
