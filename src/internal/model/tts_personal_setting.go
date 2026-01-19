package model

type TTSPersonalSetting struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	SpeakerID   string `gorm:"primaryKey;size:64"`
	SpeakerSeed int64  `gorm:"primaryKey;"`
	CreatedAt   int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt   int64  `gorm:"autoUpdateTime:nano"`
	MemberID    string `gorm:"size:64;index;not null;"`
	Member      Member `gorm:"foreignKey:MemberID;references:DiscordID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
