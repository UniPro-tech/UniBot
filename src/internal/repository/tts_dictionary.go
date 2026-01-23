package repository

import (
	"errors"

	"gorm.io/gorm"

	"unibot/internal/model"
)

type TTSDictionaryRepository struct {
	db *gorm.DB
}

// 新しいTTSDictionaryリポジトリを作成する
func NewTTSDictionaryRepository(db *gorm.DB) *TTSDictionaryRepository {
	return &TTSDictionaryRepository{db: db}
}

/* ---------------------- CRUD Methods ------------------ */

// TTS辞書エントリを作成する関数
func (r *TTSDictionaryRepository) Create(dictionary *model.TTSDictionary) error {
	return r.db.Create(dictionary).Error
}

// TTS辞書エントリをIDで取得する関数
func (r *TTSDictionaryRepository) GetByID(id uint) (*model.TTSDictionary, error) {
	var dictionary model.TTSDictionary
	err := r.db.First(&dictionary, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dictionary, err
}

// TTS辞書エントリをギルドIDとユーザーIDと単語で取得する関数
func (r *TTSDictionaryRepository) GetByGuildUserWord(guildID, userID, word string) (*model.TTSDictionary, error) {
	var dictionary model.TTSDictionary
	err := r.db.First(&dictionary, "guild_id = ? AND user_id = ? AND word = ?", guildID, userID, word).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dictionary, err
}

// TTS辞書エントリをギルドIDと単語で取得する関数
func (r *TTSDictionaryRepository) GetByGuildWord(guildID, word string) (*model.TTSDictionary, error) {
	var dictionary model.TTSDictionary
	err := r.db.First(&dictionary, "guild_id = ? AND word = ?", guildID, word).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dictionary, err
}

// TTS辞書エントリをギルドIDで全て取得する関数
func (r *TTSDictionaryRepository) ListByGuild(guildID string) ([]*model.TTSDictionary, error) {
	var dictionaries []*model.TTSDictionary
	err := r.db.Where("guild_id = ?", guildID).Order("created_at ASC").Find(&dictionaries).Error
	return dictionaries, err
}

// TTS辞書エントリをギルドIDとユーザーIDで取得する関数
func (r *TTSDictionaryRepository) ListByGuildUser(guildID, userID string) ([]*model.TTSDictionary, error) {
	var dictionaries []*model.TTSDictionary
	err := r.db.Where("guild_id = ? AND user_id = ?", guildID, userID).Order("created_at ASC").Find(&dictionaries).Error
	return dictionaries, err
}

// 全てのTTS辞書エントリを取得する関数
func (r *TTSDictionaryRepository) List() ([]*model.TTSDictionary, error) {
	var dictionaries []*model.TTSDictionary
	err := r.db.Find(&dictionaries).Error
	return dictionaries, err
}

// TTS辞書エントリを更新する関数
func (r *TTSDictionaryRepository) Update(dictionary *model.TTSDictionary) error {
	return r.db.Save(dictionary).Error
}

// TTS辞書エントリをIDで削除する関数
func (r *TTSDictionaryRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.TTSDictionary{}, id).Error
}

// TTS辞書エントリをギルドIDで削除する関数
func (r *TTSDictionaryRepository) DeleteByGuild(guildID string) error {
	return r.db.Delete(&model.TTSDictionary{}, "guild_id = ?", guildID).Error
}

// TTS辞書エントリをギルドIDとユーザーIDで削除する関数
func (r *TTSDictionaryRepository) DeleteByGuildUser(guildID, userID string) error {
	return r.db.Delete(&model.TTSDictionary{}, "guild_id = ? AND user_id = ?", guildID, userID).Error
}

// TTS辞書エントリをギルドIDと単語で削除する関数
func (r *TTSDictionaryRepository) DeleteByGuildWord(guildID, word string) error {
	return r.db.Delete(&model.TTSDictionary{}, "guild_id = ? AND word = ?", guildID, word).Error
}
