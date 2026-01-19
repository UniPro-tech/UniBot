package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type MemberRepository struct {
	db *gorm.DB
}

func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

func (r *MemberRepository) Create(discordUserID string) error {
	member := model.Member{
		DiscordUserID: discordUserID,
	}
	return r.db.Create(&member).Error
}

func (r *MemberRepository) Get(discordUserID string) (*model.Member, error) {
	var member model.Member
	err := r.db.First(&member, "discord_user_id = ?", discordUserID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &member, err
}
