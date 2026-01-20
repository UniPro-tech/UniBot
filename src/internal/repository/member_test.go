package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func TestMemberGetExists(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMemberRepository(db)

	db.Create(&model.Member{DiscordID: "123456789"})

	member, err := repo.Get("123456789")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "123456789", member.DiscordID)
}

func TestMemberGetNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMemberRepository(db)

	member, err := repo.Get("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, member)
}

func TestMemberGetWithGuilds(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMemberRepository(db)

	member := &model.Member{DiscordID: "memberID"}
	guild1 := &model.Guild{DiscordID: "guild1"}
	guild2 := &model.Guild{DiscordID: "guild2"}

	db.Create(member)
	db.Create(guild1)
	db.Create(guild2)
	db.Model(member).Association("Guilds").Append(guild1, guild2)

	retrieved, err := repo.Get("memberID")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Len(t, retrieved.Guilds, 2)
}

func TestMemberGetPreloadsGuilds(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMemberRepository(db)

	member := &model.Member{DiscordID: "user123"}
	guild := &model.Guild{DiscordID: "guild456"}

	db.Create(member)
	db.Create(guild)
	db.Model(member).Association("Guilds").Append(guild)

	retrieved, err := repo.Get("user123")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, 1, len(retrieved.Guilds))
	assert.Equal(t, "guild456", retrieved.Guilds[0].DiscordID)
}

func TestMemberGetMultipleMembers(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewMemberRepository(db)

	db.Create(&model.Member{DiscordID: "user1"})
	db.Create(&model.Member{DiscordID: "user2"})

	member, err := repo.Get("user1")
	assert.NoError(t, err)
	assert.NotNil(t, member)
	assert.Equal(t, "user1", member.DiscordID)
}
