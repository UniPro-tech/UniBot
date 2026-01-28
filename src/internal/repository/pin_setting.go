package repository

import (
	"errors"

	"unibot/internal/model"

	"gorm.io/gorm"
)

type PinSettingRepository struct {
	db *gorm.DB
}

// 新しいPinSettingRepositoryを作成します。
func NewPinSettingRepository(db *gorm.DB) *PinSettingRepository {
	return &PinSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// PinSettingを作成する関数
func (r *PinSettingRepository) Create(pinSetting *model.PinSetting) error {
	return r.db.Create(pinSetting).Error
}

// IDでPinSettingを取得する関数
func (r *PinSettingRepository) GetByID(id string) (*model.PinSetting, error) {
	var pinSetting model.PinSetting
	result := r.db.First(&pinSetting, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pinSetting, result.Error
}

// GuildIDでPinSettingを取得する関数
func (r *PinSettingRepository) GetByGuildID(guildID string) ([]*model.PinSetting, error) {
	var pinSettings []*model.PinSetting
	result := r.db.Where("guild_id = ?", guildID).Find(&pinSettings)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return pinSettings, result.Error
}

// ChannelIDでPinSettingを取得する関数
func (r *PinSettingRepository) GetByChannelID(channelID string) ([]*model.PinSetting, error) {
	var pinSettings []*model.PinSetting
	result := r.db.Where("channel_id = ?", channelID).Find(&pinSettings)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return pinSettings, result.Error
}

// 全てのPinSettingを取得する関数
func (r *PinSettingRepository) List() ([]*model.PinSetting, error) {
	var pinSettings []*model.PinSetting
	err := r.db.Find(&pinSettings).Error
	return pinSettings, err
}

// PinSettingを更新する関数
func (r *PinSettingRepository) Update(pinSetting *model.PinSetting) error {
	return r.db.Save(pinSetting).Error
}

// IDでPinSettingを削除する関数
func (r *PinSettingRepository) Delete(id string) error {
	return r.db.Delete(&model.PinSetting{}, "id = ?", id).Error
}

// ChannelIDでPinSettingを削除する関数
func (r *PinSettingRepository) DeleteByChannelID(channelID string) error {
	return r.db.Delete(&model.PinSetting{}, "channel_id = ?", channelID).Error
}
