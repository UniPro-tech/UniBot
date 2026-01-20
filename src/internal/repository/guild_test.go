package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"),
		&gorm.Config{
			//			Logger: logger.Default.LogMode(logger.Info),
		})
	require.NoError(t, err)

	err = db.AutoMigrate(
		&model.Guild{},
		&model.Member{},
		&model.AuditLogSetting{},
		&model.BotSystemSetting{},
		&model.PinSetting{},
		&model.RSSSetting{},
		&model.ScheduleSetting{},
		&model.TTSConnection{},
		&model.TTSPersonalSetting{},
	)
	require.NoError(t, err)

	return db
}

func TestGuildCreate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	err := repo.Create("123456789")
	assert.NoError(t, err)

	var guild model.Guild
	result := db.First(&guild, "discord_id = ?", "123456789")
	assert.NoError(t, result.Error)
	assert.Equal(t, "123456789", guild.DiscordID)
}

func TestGuildGet(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	db.Create(&model.Guild{DiscordID: "123456789"})

	guild, err := repo.Get("123456789")
	assert.NoError(t, err)
	assert.NotNil(t, guild)
	assert.Equal(t, "123456789", guild.DiscordID)
}

func TestGuildGetNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	guild, err := repo.Get("nonexistent")
	assert.NoError(t, err)
	assert.Nil(t, guild)
}

func TestGuildList(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	db.Create(&model.Guild{DiscordID: "111111111"})
	db.Create(&model.Guild{DiscordID: "222222222"})

	guilds, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, guilds, 2)
}

func TestGuildUpdate(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	guild := &model.Guild{DiscordID: "123456789"}
	db.Create(guild)

	guild.DiscordID = "987654321"
	err := repo.Update(guild)
	assert.NoError(t, err)

	updated, _ := repo.Get("987654321")
	assert.NotNil(t, updated)
}

func TestGuildDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	db.Create(&model.Guild{DiscordID: "123456789"})

	err := repo.Delete("123456789")
	assert.NoError(t, err)

	guild, _ := repo.Get("123456789")
	assert.Nil(t, guild)
}

func TestGuildAddMember(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	db.Create(&model.Guild{DiscordID: "guildID"})
	db.Create(&model.Member{DiscordID: "memberID"})

	err := repo.AddMember("guildID", "memberID")
	assert.NoError(t, err)

	guild, _ := repo.Get("guildID")
	assert.Len(t, guild.Members, 1)
}

func TestGuildRemoveMember(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewGuildRepository(db)

	guild := &model.Guild{DiscordID: "guildID"}
	member := &model.Member{DiscordID: "memberID"}
	db.Create(guild)
	db.Create(member)
	db.Model(guild).Association("Members").Append(member)

	err := repo.RemoveMember("guildID", "memberID")
	assert.NoError(t, err)

	updated, _ := repo.Get("guildID")
	assert.Len(t, updated.Members, 0)
}
