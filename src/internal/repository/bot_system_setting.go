package repository

import (
	"unibot/internal/model"

	"gorm.io/gorm"
)

type BotSystemSettingRepository struct {
	db *gorm.DB
}

// 新しいBotSystemSettingRepositoryを作成します。
func NewBotSystemSettingRepository(db *gorm.DB) *BotSystemSettingRepository {
	return &BotSystemSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// BotSystemSettingを作成する関数
func (r *BotSystemSettingRepository) Create(botSystemSetting *model.BotSystemSetting) error {
	return r.db.Create(botSystemSetting).Error
}

// IDでBotSystemSettingを取得する関数
func (r *BotSystemSettingRepository) GetByID(id string) (*model.BotSystemSetting, error) {
	var botSystemSetting model.BotSystemSetting
	result := r.db.Where("id = ?", id).Limit(1).Find(&botSystemSetting)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	return &botSystemSetting, nil
}

// 全てのBotSystemSettingを取得する関数
func (r *BotSystemSettingRepository) List() ([]*model.BotSystemSetting, error) {
	var botSystemSettings []*model.BotSystemSetting
	err := r.db.Find(&botSystemSettings).Error
	return botSystemSettings, err
}

// BotSystemSettingを更新する関数
func (r *BotSystemSettingRepository) Update(botSystemSetting *model.BotSystemSetting) error {
	return r.db.Save(botSystemSetting).Error
}

// IDでBotSystemSettingを削除する関数
func (r *BotSystemSettingRepository) Delete(id string) error {
	return r.db.Delete(&model.BotSystemSetting{}, "id = ?", id).Error
}
