package model

type Member struct {
	DiscordID           string              `gorm:"primaryKey;size:255"`
	Guilds              []*Guild            `gorm:"many2many:members_guilds;"`
	TTSPersonalSettings *TTSPersonalSetting `gorm:"foreignKey:AuthorID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduleSetting     []*ScheduleSetting  `gorm:"foreignKey:ID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
