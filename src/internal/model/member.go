package model

type Member struct {
	DiscordID          string              `gorm:"primaryKey;size:255"`
	Guilds             []*Guild            `gorm:"many2many:members_guilds;"`
	TTSPersonalSetting *TTSPersonalSetting `gorm:"foreignKey:MemberID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ScheduleSetting    []*ScheduleSetting  `gorm:"foreignKey:AuthorID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
