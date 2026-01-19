package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type GuildRepository struct {
	db *gorm.DB
}

// リポジトリのインスタンスを作成する関数
func NewGuildRepository(db *gorm.DB) *GuildRepository {
	return &GuildRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// DiscordのguildIDからギルドを作成する関数
func (r *GuildRepository) Create(DiscordID string) error {
	guild := model.Guild{DiscordID: DiscordID}
	return r.db.FirstOrCreate(&guild).Error
}

// DiscordのguildIDからギルドを取得する関数
func (r *GuildRepository) Get(DiscordID string) (*model.Guild, error) {
	var guild model.Guild
	err := r.db.Preload("Members").
		First(&guild, "discord_id = ?", DiscordID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &guild, err
}

// 全てのギルドを取得する関数
func (r *GuildRepository) List() ([]*model.Guild, error) {
	var guilds []*model.Guild
	err := r.db.Find(&guilds).Error
	return guilds, err
}

// ギルド情報を更新する関数
func (r *GuildRepository) Update(guild *model.Guild) error {
	return r.db.Save(guild).Error
}

// DiscordのguildIDからギルドを削除する関数
func (r *GuildRepository) Delete(DiscordID string) error {
	return r.db.Delete(&model.Guild{}, "discord_id = ?", DiscordID).Error
}

/* ---------------------- Relational Methods ------------------ */

// ギルドにメンバーを追加する関数
func (r *GuildRepository) AddMember(
	guildID string,
	memberID string,
) error {
	var guild model.Guild
	if err := r.db.First(&guild, "discord_id = ?", guildID).Error; err != nil {
		return err
	}

	var member model.Member
	if err := r.db.First(&member, "discord_id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&guild).Association("Members").Append(&member)
}

// メンバーからギルドを削除する関数
func (r *GuildRepository) RemoveMember(
	guildID string,
	memberID string,
) error {
	var guild model.Guild
	if err := r.db.Preload("Members").
		First(&guild, "discord_id = ?", guildID).Error; err != nil {
		return err
	}

	var member model.Member
	if err := r.db.First(&member, "discord_id = ?", memberID).Error; err != nil {
		return err
	}

	return r.db.Model(&guild).Association("Members").Delete(&member)
}
