package model

import "time"

// RolePanel セレクトメニュー式ロールパネルのモデル
type RolePanel struct {
	ID          uint               `gorm:"primaryKey;autoIncrement"`
	GuildID     string             `gorm:"size:255;not null;index"`
	ChannelID   string             `gorm:"size:255;not null"`
	MessageID   string             `gorm:"size:255;not null;uniqueIndex"`
	Title       string             `gorm:"size:255;not null"`
	Description string             `gorm:"size:1000"`
	Options     []*RolePanelOption `gorm:"foreignKey:RolePanelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RolePanelOption ロールパネルの選択肢
type RolePanelOption struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	RolePanelID uint   `gorm:"not null;index"`
	OptionKey   string `gorm:"size:64;not null;uniqueIndex"`
	Label       string `gorm:"size:100;not null"`
	Description string `gorm:"size:100"`
	Emoji       string `gorm:"size:255"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
