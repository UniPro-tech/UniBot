package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type RolePanelRepository struct {
	db *gorm.DB
}

// 新しいRolePanelリポジトリを作成する
func NewRolePanelRepository(db *gorm.DB) *RolePanelRepository {
	return &RolePanelRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// ロールパネルを作成する関数
func (r *RolePanelRepository) Create(panel *model.RolePanel) error {
	return r.db.Create(panel).Error
}

// ロールパネルをIDで取得する関数
func (r *RolePanelRepository) GetByID(id uint) (*model.RolePanel, error) {
	var panel model.RolePanel
	err := r.db.Preload("Options").First(&panel, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &panel, err
}

// ロールパネルをメッセージIDで取得する関数
func (r *RolePanelRepository) GetByMessageID(messageID string) (*model.RolePanel, error) {
	var panel model.RolePanel
	err := r.db.Preload("Options").First(&panel, "message_id = ?", messageID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &panel, err
}

// ロールパネルをギルドIDで全て取得する関数
func (r *RolePanelRepository) ListByGuild(guildID string) ([]*model.RolePanel, error) {
	var panels []*model.RolePanel
	err := r.db.Preload("Options").Where("guild_id = ?", guildID).Order("created_at ASC").Find(&panels).Error
	return panels, err
}

// ロールパネルをチャンネルIDで全て取得する関数
func (r *RolePanelRepository) ListByChannel(channelID string) ([]*model.RolePanel, error) {
	var panels []*model.RolePanel
	err := r.db.Preload("Options").Where("channel_id = ?", channelID).Order("created_at ASC").Find(&panels).Error
	return panels, err
}

// 全てのロールパネルを取得する関数
func (r *RolePanelRepository) List() ([]*model.RolePanel, error) {
	var panels []*model.RolePanel
	err := r.db.Preload("Options").Find(&panels).Error
	return panels, err
}

// ロールパネルを更新する関数
func (r *RolePanelRepository) Update(panel *model.RolePanel) error {
	return r.db.Save(panel).Error
}

// ロールパネルをIDで削除する関数
func (r *RolePanelRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.RolePanel{}, id).Error
}

// ロールパネルをメッセージIDで削除する関数
func (r *RolePanelRepository) DeleteByMessageID(messageID string) error {
	return r.db.Delete(&model.RolePanel{}, "message_id = ?", messageID).Error
}

// ロールパネルをギルドIDで削除する関数
func (r *RolePanelRepository) DeleteByGuild(guildID string) error {
	return r.db.Delete(&model.RolePanel{}, "guild_id = ?", guildID).Error
}

/* ---------------------- Option Methods ------------------ */

// ロールパネルにオプションを追加する関数
func (r *RolePanelRepository) AddOption(option *model.RolePanelOption) error {
	return r.db.Create(option).Error
}

// オプションをIDで取得する関数
func (r *RolePanelRepository) GetOptionByID(id uint) (*model.RolePanelOption, error) {
	var option model.RolePanelOption
	err := r.db.First(&option, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &option, err
}

// オプションをロールパネルIDで取得する関数
func (r *RolePanelRepository) ListOptionsByPanelID(panelID uint) ([]*model.RolePanelOption, error) {
	var options []*model.RolePanelOption
	err := r.db.Where("role_panel_id = ?", panelID).Order("created_at ASC").Find(&options).Error
	return options, err
}

// オプションを更新する関数
func (r *RolePanelRepository) UpdateOption(option *model.RolePanelOption) error {
	return r.db.Save(option).Error
}

// オプションをIDで削除する関数
func (r *RolePanelRepository) DeleteOptionByID(id uint) error {
	return r.db.Delete(&model.RolePanelOption{}, id).Error
}

// オプションをロールパネルIDで全て削除する関数
func (r *RolePanelRepository) DeleteOptionsByPanelID(panelID uint) error {
	return r.db.Delete(&model.RolePanelOption{}, "role_panel_id = ?", panelID).Error
}
