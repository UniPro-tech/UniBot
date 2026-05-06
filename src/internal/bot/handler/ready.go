package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"unibot/internal"
	"unibot/internal/bot/command/admin/maintenance"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Ready(ctx *internal.BotContext, e *events.Ready) {
	log.Println("Bot is ready 🚀")
	log.Printf("Logged in as: %v#%v", e.User.Username, e.User.Discriminator)

	// Get DB Connection
	database := ctx.DB

	// Get status
	settingsRepo := repository.NewBotSystemSettingRepository(database)
	botStatusRecord, err := settingsRepo.GetByID("status")
	if err != nil {
		return
	} else {
		if botStatusRecord == nil {
			serverCounts := len(e.Guilds)
			statusData := maintenance.StatusData{
				Text:         fmt.Sprintf("Serving %d servers | /help", serverCounts),
				Type:         discord.ActivityTypeGame,
				OnlineStatus: discord.OnlineStatusOnline,
			}
			if err := maintenance.SetBotStatus(e.Client(), statusData); err != nil {
				return
			}
			newSetting := &model.BotSystemSetting{ID: "status"}
			if err := newSetting.Value.Set(statusData); err != nil {
				return
			}
			if err := settingsRepo.Create(newSetting); err != nil {
				log.Printf("Error creating default bot status: %v", err)
			}
		} else {
			botStatus := botStatusRecord.Value

			// jsonからデータを引き出す処理を追加
			var statusData maintenance.StatusData
			err = json.Unmarshal(botStatus.Bytes, &statusData)

			if err != nil {
				return
			} else {
				// Set Bot Status
				err = maintenance.SetBotStatus(e.Client(), statusData)
				if err != nil {
					log.Printf("Error setting bot status: %v", err)
					return
				}
			}
		}
	}
}
