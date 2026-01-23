package model

// TTS辞書エントリ
type TTSDictionary struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	GuildID       string `gorm:"size:64;not null;index"`
	UserID        string `gorm:"size:64;not null;index"`
	Word          string `gorm:"size:255;not null"`
	Definition    string `gorm:"size:255;not null"`
	CaseSensitive bool   `gorm:"default:false"`
	CreatedAt     int64  `gorm:"autoCreateTime:nano"`
	UpdatedAt     int64  `gorm:"autoUpdateTime:nano"`
}
