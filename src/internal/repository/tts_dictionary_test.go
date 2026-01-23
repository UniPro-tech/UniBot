package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func setupDictTestDB(t *testing.T) *repository.TTSDictionaryRepository {
	db := setupTestDB(t)
	db.AutoMigrate(&model.TTSDictionary{})
	return repository.NewTTSDictionaryRepository(db)
}

func TestTTSDictionaryCreate(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user456",
		Word:          "テスト",
		Definition:    "てすと",
		CaseSensitive: false,
	}

	err := repo.Create(entry)
	assert.NoError(t, err)
	assert.NotZero(t, entry.ID)
}

func TestTTSDictionaryGetByID(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:    "guild123",
		UserID:     "user456",
		Word:       "テスト",
		Definition: "てすと",
	}
	repo.Create(entry)

	retrieved, err := repo.GetByID(entry.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "テスト", retrieved.Word)
}

func TestTTSDictionaryGetByIDNotFound(t *testing.T) {
	repo := setupDictTestDB(t)

	retrieved, err := repo.GetByID(9999)
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestTTSDictionaryGetByGuildWord(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:    "guild123",
		UserID:     "user456",
		Word:       "テスト",
		Definition: "てすと",
	}
	repo.Create(entry)

	retrieved, err := repo.GetByGuildWord("guild123", "テスト")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "てすと", retrieved.Definition)
}

func TestTTSDictionaryListByGuild(t *testing.T) {
	repo := setupDictTestDB(t)

	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語1", Definition: "読み1"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user2", Word: "単語2", Definition: "読み2"})
	repo.Create(&model.TTSDictionary{GuildID: "guild456", UserID: "user1", Word: "単語3", Definition: "読み3"})

	entries, err := repo.ListByGuild("guild123")
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
}

func TestTTSDictionaryListByGuildUser(t *testing.T) {
	repo := setupDictTestDB(t)

	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語1", Definition: "読み1"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user2", Word: "単語2", Definition: "読み2"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語3", Definition: "読み3"})

	entries, err := repo.ListByGuildUser("guild123", "user1")
	assert.NoError(t, err)
	assert.Len(t, entries, 2)
}

func TestTTSDictionaryUpdate(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:    "guild123",
		UserID:     "user456",
		Word:       "テスト",
		Definition: "てすと",
	}
	repo.Create(entry)

	entry.Definition = "更新後の読み"
	err := repo.Update(entry)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByID(entry.ID)
	assert.Equal(t, "更新後の読み", retrieved.Definition)
}

func TestTTSDictionaryDeleteByID(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:    "guild123",
		UserID:     "user456",
		Word:       "テスト",
		Definition: "てすと",
	}
	repo.Create(entry)

	err := repo.DeleteByID(entry.ID)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByID(entry.ID)
	assert.Nil(t, retrieved)
}

func TestTTSDictionaryDeleteByGuild(t *testing.T) {
	repo := setupDictTestDB(t)

	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語1", Definition: "読み1"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user2", Word: "単語2", Definition: "読み2"})

	err := repo.DeleteByGuild("guild123")
	assert.NoError(t, err)

	entries, _ := repo.ListByGuild("guild123")
	assert.Len(t, entries, 0)
}

func TestTTSDictionaryDeleteByGuildUser(t *testing.T) {
	repo := setupDictTestDB(t)

	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語1", Definition: "読み1"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user1", Word: "単語2", Definition: "読み2"})
	repo.Create(&model.TTSDictionary{GuildID: "guild123", UserID: "user2", Word: "単語3", Definition: "読み3"})

	err := repo.DeleteByGuildUser("guild123", "user1")
	assert.NoError(t, err)

	entries, _ := repo.ListByGuild("guild123")
	assert.Len(t, entries, 1)
	assert.Equal(t, "user2", entries[0].UserID)
}

func TestTTSDictionaryCaseSensitive(t *testing.T) {
	repo := setupDictTestDB(t)

	entry := &model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user456",
		Word:          "Test",
		Definition:    "テスト",
		CaseSensitive: true,
	}
	repo.Create(entry)

	retrieved, _ := repo.GetByID(entry.ID)
	assert.True(t, retrieved.CaseSensitive)
}
