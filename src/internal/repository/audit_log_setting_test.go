package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func TestAuditLogSettingCreateSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting := &model.AuditLogSetting{
		GuildID:        "guild123",
		AlartChannelID: "channel456",
	}

	err := repo.Create(setting)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByGuildID("guild123")
	assert.NotNil(t, retrieved)
	assert.Equal(t, "guild123", retrieved.GuildID)
	assert.Equal(t, "channel456", retrieved.AlartChannelID)
}

func TestAuditLogSettingGetByGuildIDExists(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting := &model.AuditLogSetting{
		GuildID:        "guild789",
		AlartChannelID: "channel999",
	}
	db.Create(setting)

	retrieved, err := repo.GetByGuildID("guild789")
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, "guild789", retrieved.GuildID)
	assert.Equal(t, "channel999", retrieved.AlartChannelID)
}

func TestAuditLogSettingGetByGuildIDNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	retrieved, err := repo.GetByGuildID("nonexistent_guild")
	assert.NoError(t, err)
	assert.Nil(t, retrieved)
}

func TestAuditLogSettingUpdateSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	original := &model.AuditLogSetting{
		GuildID:        "guild111",
		AlartChannelID: "channel111",
	}
	db.Create(original)

	original.AlartChannelID = "channel_updated"
	err := repo.Update(original)
	assert.NoError(t, err)

	updated, _ := repo.GetByGuildID("guild111")
	assert.NotNil(t, updated)
	assert.Equal(t, "channel_updated", updated.AlartChannelID)
}

func TestAuditLogSettingDeleteByGuildIDSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting := &model.AuditLogSetting{
		GuildID:        "guild_to_delete",
		AlartChannelID: "channel_temp",
	}
	db.Create(setting)

	err := repo.DeleteByGuildID("guild_to_delete")
	assert.NoError(t, err)

	retrieved, _ := repo.GetByGuildID("guild_to_delete")
	assert.Nil(t, retrieved)
}

func TestAuditLogSettingMultipleGuilds(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting1 := &model.AuditLogSetting{GuildID: "guild1", AlartChannelID: "ch1"}
	setting2 := &model.AuditLogSetting{GuildID: "guild2", AlartChannelID: "ch2"}
	repo.Create(setting1)
	repo.Create(setting2)

	retrieved1, _ := repo.GetByGuildID("guild1")
	retrieved2, _ := repo.GetByGuildID("guild2")

	assert.Equal(t, "guild1", retrieved1.GuildID)
	assert.Equal(t, "guild2", retrieved2.GuildID)
	assert.NotEqual(t, retrieved1.AlartChannelID, retrieved2.AlartChannelID)
}

func TestAuditLogSettingDeletePreservesOthers(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting1 := &model.AuditLogSetting{GuildID: "guild_keep", AlartChannelID: "ch_keep"}
	setting2 := &model.AuditLogSetting{GuildID: "guild_delete", AlartChannelID: "ch_delete"}
	repo.Create(setting1)
	repo.Create(setting2)

	repo.DeleteByGuildID("guild_delete")

	kept, _ := repo.GetByGuildID("guild_keep")
	deleted, _ := repo.GetByGuildID("guild_delete")

	assert.NotNil(t, kept)
	assert.Nil(t, deleted)
}

func TestAuditLogSettingUpdateMultipleFields(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewAuditLogSettingRepository(db)

	setting := &model.AuditLogSetting{
		GuildID:        "guild_multi",
		AlartChannelID: "ch_original",
	}
	db.Create(setting)

	setting.AlartChannelID = "ch_updated"
	repo.Update(setting)

	updated, err := repo.GetByGuildID("guild_multi")
	assert.NoError(t, err)
	assert.Equal(t, "ch_updated", updated.AlartChannelID)
}
