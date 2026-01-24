package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func TestDictionaryCache_IsolatedByDB(t *testing.T) {
	cache := NewDictionaryCache(time.Minute)
	db1 := setupTestDB(t)
	db2 := setupTestDB(t)

	repo1 := repository.NewTTSDictionaryRepository(db1)
	repo2 := repository.NewTTSDictionaryRepository(db2)

	err := repo1.Create(&model.TTSDictionary{
		GuildID:       "guildA",
		UserID:        "user1",
		Word:          "foo",
		Definition:    "bar",
		CaseSensitive: false,
	})
	require.NoError(t, err)

	err = repo2.Create(&model.TTSDictionary{
		GuildID:       "guildA",
		UserID:        "user2",
		Word:          "baz",
		Definition:    "qux",
		CaseSensitive: false,
	})
	require.NoError(t, err)

	entries1, err := cache.Get(db1, "guildA")
	require.NoError(t, err)
	entries2, err := cache.Get(db2, "guildA")
	require.NoError(t, err)

	if assert.Len(t, entries1, 1) {
		assert.Equal(t, "foo", entries1[0].Word)
	}
	if assert.Len(t, entries2, 1) {
		assert.Equal(t, "baz", entries2[0].Word)
	}
}

func TestDictionaryCache_Invalidate(t *testing.T) {
	cache := NewDictionaryCache(time.Minute)
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	entry := &model.TTSDictionary{
		GuildID:       "guild1",
		UserID:        "user1",
		Word:          "hello",
		Definition:    "hi",
		CaseSensitive: false,
	}
	err := repo.Create(entry)
	require.NoError(t, err)

	entries, err := cache.Get(db, "guild1")
	require.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "hello", entries[0].Word)
	}

	err = repo.DeleteByID(entry.ID)
	require.NoError(t, err)

	entries, err = cache.Get(db, "guild1")
	require.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "hello", entries[0].Word)
	}

	cache.Invalidate("guild1")
	entries, err = cache.Get(db, "guild1")
	require.NoError(t, err)
	assert.Len(t, entries, 0)
}

func TestDictionaryCache_TTLRefresh(t *testing.T) {
	cache := NewDictionaryCache(10 * time.Millisecond)
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	entry := &model.TTSDictionary{
		GuildID:       "guildTTL",
		UserID:        "user1",
		Word:          "A",
		Definition:    "B",
		CaseSensitive: false,
	}
	err := repo.Create(entry)
	require.NoError(t, err)

	entries, err := cache.Get(db, "guildTTL")
	require.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "A", entries[0].Word)
	}

	err = repo.DeleteByID(entry.ID)
	require.NoError(t, err)

	err = repo.Create(&model.TTSDictionary{
		GuildID:       "guildTTL",
		UserID:        "user1",
		Word:          "C",
		Definition:    "D",
		CaseSensitive: false,
	})
	require.NoError(t, err)

	entries, err = cache.Get(db, "guildTTL")
	require.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "A", entries[0].Word)
	}

	time.Sleep(20 * time.Millisecond)

	entries, err = cache.Get(db, "guildTTL")
	require.NoError(t, err)
	if assert.Len(t, entries, 1) {
		assert.Equal(t, "C", entries[0].Word)
	}
}
