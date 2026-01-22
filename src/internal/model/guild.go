package model

type Guild struct {
	DiscordID        string             `gorm:"primaryKey;size:255"`
	Members          []*Member          `gorm:"many2many:members_guilds;"`
	RSSSettings      []*RSSSetting      `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduleSetting  []*ScheduleSetting `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AuditLogSettings []*AuditLogSetting `gorm:"foreignKey:GuildID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
