package model

type TTSPersonalSetting struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	SpeakerID    string `gorm:"size:16;not null;"`
	SpeakerPitch int64
	CreatedAt    int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt    int64  `gorm:"autoUpdateTime:nano"`
	MemberID     string `gorm:"size:64;index;not null;"`
}
