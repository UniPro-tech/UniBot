package model

type TTSPersonalSetting struct {
	MemberID    string `gorm:"primaryKey;size:255;uniqueIndex:idx_author_speaker"`
	SpeakerID   string `gorm:"primaryKey;size:64"`
	SpeakerSeed int64  `gorm:"primaryKey;"`
	CreatedAt   int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:nano"`
	Member      Member `gorm:"foreignKey:MemberID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
