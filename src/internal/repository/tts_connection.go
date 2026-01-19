package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type TTSConnectionRepository struct {
	db *gorm.DB
}

// 新しいTTSConnectionリポジトリを作成する
func NewTTSConnectionRepository(db *gorm.DB) *TTSConnectionRepository {
	return &TTSConnectionRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// TTS接続を作成する関数
func (r *TTSConnectionRepository) Create(ttsConnection *model.TTSConnection) error {
	return r.db.Create(ttsConnection).Error
}

// TTS接続をギルドIDで取得する関数
func (r *TTSConnectionRepository) GetByGuildID(guildID string) (*model.TTSConnection, error) {
	var ttsConnection model.TTSConnection
	err := r.db.First(&ttsConnection, "guild_id = ?", guildID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ttsConnection, err
}

// TTS接続をチャンネルIDで取得する関数
func (r *TTSConnectionRepository) GetByChannelID(channelID string) (*model.TTSConnection, error) {
	var ttsConnection model.TTSConnection
	err := r.db.First(&ttsConnection, "channel_id = ?", channelID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ttsConnection, err
}

// 全てのTTS接続を取得する関数
func (r *TTSConnectionRepository) List() ([]*model.TTSConnection, error) {
	var ttsConnections []*model.TTSConnection
	err := r.db.Find(&ttsConnections).Error
	return ttsConnections, err
}

// TTS接続を更新する関数
func (r *TTSConnectionRepository) Update(ttsConnection *model.TTSConnection) error {
	return r.db.Save(ttsConnection).Error
}

// TTS接続をギルドIDで削除する関数
func (r *TTSConnectionRepository) DeleteByGuildID(guildID string) error {
	return r.db.Delete(&model.TTSConnection{}, "guild_id = ?", guildID).Error
}
