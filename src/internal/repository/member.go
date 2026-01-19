package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type MemberRepository struct {
	db *gorm.DB
}

// リポジトリのインスタンスを作成する関数
func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// DiscordのuserIDからメンバーを作成する関数
func (r *MemberRepository) Create(DiscordID string) error {
	member := model.Member{DiscordID: DiscordID}
	return r.db.FirstOrCreate(&member).Error
}

// DiscordのuserIDからメンバーを取得する関数
func (r *MemberRepository) Get(DiscordID string) (*model.Member, error) {
	var member model.Member
	err := r.db.Preload("Guilds").
		First(&member, "discord_id = ?", DiscordID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &member, err
}

// 全てのメンバーを取得する関数
func (r *MemberRepository) List() ([]*model.Member, error) {
	var members []*model.Member
	err := r.db.Find(&members).Error
	return members, err
}

// メンバー情報を更新する関数
func (r *MemberRepository) Update(member *model.Member) error {
	return r.db.Save(member).Error
}

// DiscordのuserIDからメンバーを削除する関数
func (r *MemberRepository) Delete(DiscordID string) error {
	return r.db.Delete(&model.Member{}, "discord_id = ?", DiscordID).Error
}

/* ---------------------- Relational Methods ------------------ */

// メンバーにギルドを追加する関数
func (r *MemberRepository) AddGuild(
	memberID string,
	guildID string,
) error {
	var member model.Member
	if err := r.db.Preload("Guilds").
		First(&member, "discord_id = ?", memberID).Error; err != nil {
		return err
	}

	var guild model.Guild
	if err := r.db.First(&guild, "discord_id = ?", guildID).Error; err != nil {
		return err
	}

	return r.db.Model(&member).Association("Guilds").Append(&guild)
}

// メンバーからギルドを削除する関数
func (r *MemberRepository) RemoveGuild(
	memberID string,
	guildID string,
) error {
	var member model.Member
	if err := r.db.Preload("Guilds").
		First(&member, "discord_id = ?", memberID).Error; err != nil {
		return err
	}

	var guild model.Guild
	if err := r.db.First(&guild, "discord_id = ?", guildID).Error; err != nil {
		return err
	}

	return r.db.Model(&member).Association("Guilds").Delete(&guild)
}
