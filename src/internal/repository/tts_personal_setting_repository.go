package repository

import (
	"errors"

	"unibot/internal/model"

	"gorm.io/gorm"
)

type TTSPersonalSettingRepository struct {
	db *gorm.DB
}

// 新しいTTSPersonalSettingリポジトリを作成する
func NewTTSPersonalSettingRepository(db *gorm.DB) *TTSPersonalSettingRepository {
	return &TTSPersonalSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// TTSPersonalSettingを作成する関数
func (r *TTSPersonalSettingRepository) Create(ttsPersonalSetting *model.TTSPersonalSetting) error {
	return r.db.Create(ttsPersonalSetting).Error
}

// TTSPersonalSettingをMemberIDで取得する関数
func (r *TTSPersonalSettingRepository) GetByMember(memberID string) (*model.TTSPersonalSetting, error) {
	var ttsPersonalSetting model.TTSPersonalSetting
	err := r.db.First(&ttsPersonalSetting, "member_id = ?", memberID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ttsPersonalSetting, err
}

// 全てのTTSPersonalSettingを取得する関数
func (r *TTSPersonalSettingRepository) List() ([]*model.TTSPersonalSetting, error) {
	var ttsPersonalSettings []*model.TTSPersonalSetting
	err := r.db.Find(&ttsPersonalSettings).Error
	return ttsPersonalSettings, err
}

// TTSPersonalSettingを更新する関数
func (r *TTSPersonalSettingRepository) Update(ttsPersonalSetting *model.TTSPersonalSetting) error {
	return r.db.Save(ttsPersonalSetting).Error
}

// TTSPersonalSettingをMemberIDで削除する関数
func (r *TTSPersonalSettingRepository) DeleteByMember(memberID string) error {
	return r.db.Delete(&model.TTSPersonalSetting{}, "member_id = ?", memberID).Error
}
