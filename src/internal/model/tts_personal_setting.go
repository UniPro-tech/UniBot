package model

type TTSPersonalSetting struct {
	AuthorID    string `gorm:"primaryKey;size:255;uniqueIndex:idx_author_speaker"`
	SpeakerID   string `gorm:"primaryKey;size:64"`
	SpeakerSeed int64  `gorm:"primaryKey;"`
	CreatedAt   int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:nano"`
	Author      Member `gorm:"foreignKey:AuthorID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
