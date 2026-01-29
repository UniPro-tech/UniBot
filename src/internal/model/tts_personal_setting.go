package model

type TTSPersonalSetting struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	SpeakerID    string `gorm:"size:16;not null;"`
	SpeakerPitch int64
	SpeedScale   int64  `gorm:"default:100"` // 100 = `1.0倍速`（整数で管理）
	CreatedAt    int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt    int64  `gorm:"autoUpdateTime:nano"`
	MemberID     string `gorm:"size:64;index;not null;"`
}
