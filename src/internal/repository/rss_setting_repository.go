package repository

import (
	"errors"

	"unibot/internal/model"

	"gorm.io/gorm"
)

type RSSSettingRepository struct {
	db *gorm.DB
}

// 新しいRSSSettingリポジトリを作成する
func NewRSSSettingRepository(db *gorm.DB) *RSSSettingRepository {
	return &RSSSettingRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// RSS設定を作成する関数
func (r *RSSSettingRepository) Create(rssSetting *model.RSSSetting) error {
	return r.db.Create(rssSetting).Error
}

// RSS設定をIDで取得する関数
func (r *RSSSettingRepository) GetByID(id string) (*model.RSSSetting, error) {
	var rssSetting model.RSSSetting
	err := r.db.First(&rssSetting, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rssSetting, err
}

// 全てのRSS設定を取得する関数
func (r *RSSSettingRepository) List() ([]*model.RSSSetting, error) {
	var rssSettings []*model.RSSSetting
	err := r.db.Find(&rssSettings).Error
	return rssSettings, err
}

// RSS設定を更新する関数
func (r *RSSSettingRepository) Update(rssSetting *model.RSSSetting) error {
	return r.db.Save(rssSetting).Error
}

// RSS設定をIDで削除する関数
func (r *RSSSettingRepository) DeleteByID(id string) error {
	return r.db.Delete(&model.RSSSetting{}, "id = ?", id).Error
}
