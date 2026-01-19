package model

type Guild struct {
	DiscordID   string        `gorm:"primaryKey;size:255"`
	Members     []*Member     `gorm:"many2many:members_guilds;"`
	RSSSettings []*RSSSetting `gorm:"foreignKey:ChannelID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
