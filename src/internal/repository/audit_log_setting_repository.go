package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type AuditLogSettingRepository struct {
	db *gorm.DB
}

// リポジトリのインスタンスを作成する関数
func NewAuditLogSettingRepository(db *gorm.DB) *AuditLogSettingRepository {
	return &AuditLogSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// 監査ログ設定を作成する関数
func (r *AuditLogSettingRepository) Create(setting *model.AuditLogSetting) error {
	return r.db.Create(setting).Error
}

// ギルドIDから監査ログ設定を取得する関数
func (r *AuditLogSettingRepository) GetByGuildID(guildID string) (*model.AuditLogSetting, error) {
	var setting model.AuditLogSetting
	err := r.db.First(&setting, "guild_id = ?", guildID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &setting, err
}

// 監査ログ設定を更新する関数
func (r *AuditLogSettingRepository) Update(setting *model.AuditLogSetting) error {
	return r.db.Save(setting).Error
}

// ギルドIDから監査ログ設定を削除する関数
func (r *AuditLogSettingRepository) DeleteByGuildID(guildID string) error {
	return r.db.Delete(&model.AuditLogSetting{}, "guild_id = ?", guildID).Error
}
