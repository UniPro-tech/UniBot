package model

type RSSSetting struct {
	ID                           string  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	URL                          string  `gorm:"not null"`
	ChannelID                    string  `gorm:"not null"`
	WebhookURL                   string  `gorm:"not null"`
	Title                        *string `gorm:"null"`
	IsFailed                     bool    `gorm:"default:false"`
	LastItemTitleDescriptionHash *string `gorm:"null"`
	CreatedAt                    int64   `gorm:"autoCreateTime:nano"`
	UpdatedAt                    int64   `gorm:"autoUpdateTime:nano"`
	GuildID                      string  `gorm:"not null"`
	Guild                        Guild   `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
