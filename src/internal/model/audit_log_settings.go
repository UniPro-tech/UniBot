package model

type AuditLogSettings struct {
	GuildID        string `gorm:"primaryKey;size:255"`
	AlartChannelID string `gorm:"size:255"`
	AlartType      int    `gorm:"not null;default:0"`
	CreatedAt      int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt      int64  `gorm:"autoUpdateTime:nano"`
	Guild          Guild  `gorm:"foreignKey:GuildID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
