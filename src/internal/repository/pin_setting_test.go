package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func TestPinSettingUpdateSuccess(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPinSettingRepository(db)

	original := &model.PinSetting{
		ID:        "pin123",
		GuildID:   "guild1",
		ChannelID: "channel1",
	}
	db.Create(original)

	original.ChannelID = "channel2"
	err := repo.Update(original)
	assert.NoError(t, err)

	updated, _ := repo.GetByID("pin123")
	assert.NotNil(t, updated)
	assert.Equal(t, "channel2", updated.ChannelID)
}

func TestPinSettingUpdateMultipleFields(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPinSettingRepository(db)

	pinSetting := &model.PinSetting{
		ID:        "pin456",
		GuildID:   "guild1",
		ChannelID: "channel1",
	}
	db.Create(pinSetting)

	pinSetting.GuildID = "guild2"
	pinSetting.ChannelID = "channel2"
	err := repo.Update(pinSetting)
	assert.NoError(t, err)

	updated, _ := repo.GetByID("pin456")
	assert.NotNil(t, updated)
	assert.Equal(t, "guild2", updated.GuildID)
	assert.Equal(t, "channel2", updated.ChannelID)
}

func TestPinSettingUpdateSameValues(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPinSettingRepository(db)

	pinSetting := &model.PinSetting{
		ID:        "pin789",
		GuildID:   "guild1",
		ChannelID: "channel1",
	}
	db.Create(pinSetting)

	err := repo.Update(pinSetting)
	assert.NoError(t, err)

	updated, _ := repo.GetByID("pin789")
	assert.NotNil(t, updated)
	assert.Equal(t, "guild1", updated.GuildID)
	assert.Equal(t, "channel1", updated.ChannelID)
}

func TestPinSettingUpdatePreservesOtherRecords(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPinSettingRepository(db)

	pin1 := &model.PinSetting{ID: "pin1", GuildID: "guild1", ChannelID: "ch1"}
	pin2 := &model.PinSetting{ID: "pin2", GuildID: "guild2", ChannelID: "ch2"}
	db.Create(pin1)
	db.Create(pin2)

	pin1.ChannelID = "ch_updated"
	err := repo.Update(pin1)
	assert.NoError(t, err)

	updated1, _ := repo.GetByID("pin1")
	updated2, _ := repo.GetByID("pin2")
	assert.Equal(t, "ch_updated", updated1.ChannelID)
	assert.Equal(t, "ch2", updated2.ChannelID)
}

func TestPinSettingUpdateNonexistent(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewPinSettingRepository(db)

	nonexistent := &model.PinSetting{
		ID:        "nonexistent",
		GuildID:   "guild1",
		ChannelID: "channel1",
	}

	err := repo.Update(nonexistent)
	assert.NoError(t, err)

	retrieved, _ := repo.GetByID("nonexistent")
	assert.NotNil(t, retrieved)
}
