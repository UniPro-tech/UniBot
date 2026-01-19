package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type ScheduleSettingRepository struct {
	db *gorm.DB
}

// 新しいScheduleSettingリポジトリを作成する
func NewScheduleSettingRepository(db *gorm.DB) *ScheduleSettingRepository {
	return &ScheduleSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// スケジュール設定を作成する関数
func (r *ScheduleSettingRepository) Create(scheduleSetting *model.ScheduleSetting) error {
	return r.db.Create(scheduleSetting).Error
}

// スケジュール設定をIDで取得する関数
func (r *ScheduleSettingRepository) GetByID(id string) (*model.ScheduleSetting, error) {
	var scheduleSetting model.ScheduleSetting
	err := r.db.First(&scheduleSetting, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &scheduleSetting, err
}

// スケジュール設定をGuildIDで取得する関数
func (r *ScheduleSettingRepository) GetByGuildID(guildID string) ([]*model.ScheduleSetting, error) {
	var scheduleSettings []*model.ScheduleSetting
	err := r.db.Where("guild_id = ?", guildID).Find(&scheduleSettings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return scheduleSettings, err
}

// スケジュール設定をChannelIDで取得する関数
func (r *ScheduleSettingRepository) GetByChannelID(channelID string) ([]*model.ScheduleSetting, error) {
	var scheduleSettings []*model.ScheduleSetting
	err := r.db.Where("channel_id = ?", channelID).Find(&scheduleSettings).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return scheduleSettings, err
}

// 全てのスケジュール設定を取得する関数
func (r *ScheduleSettingRepository) List() ([]*model.ScheduleSetting, error) {
	var scheduleSettings []*model.ScheduleSetting
	err := r.db.Find(&scheduleSettings).Error
	return scheduleSettings, err
}

// スケジュール設定を更新する関数
func (r *ScheduleSettingRepository) Update(scheduleSetting *model.ScheduleSetting) error {
	return r.db.Save(scheduleSetting).Error
}

// スケジュール設定をIDで削除する関数
func (r *ScheduleSettingRepository) DeleteByID(id string) error {
	return r.db.Delete(&model.ScheduleSetting{}, "id = ?", id).Error
}
