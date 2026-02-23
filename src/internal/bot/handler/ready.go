package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"unibot/internal"
	"unibot/internal/bot/command/admin/maintenance"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func Ready(ctx *internal.BotContext) func(s *discordgo.Session, r *discordgo.Ready) {
	return func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is ready 🚀")
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)

		// Get DB Connection
		database := ctx.DB

		// Get status
		settingsRepo := repository.NewBotSystemSettingRepository(database)
		botStatusRecord, err := settingsRepo.GetByID("status")
		if err != nil {
			log.Printf("Error getting bot status: %v", err)
			return
		} else {
			if botStatusRecord == nil {
				log.Println("Bot status not found in the database.")
				serverCounts := s.State.Guilds
				statusData := maintenance.StatusData{
					Text:         "Serving " + fmt.Sprintf("%d", len(serverCounts)) + " servers | /help",
					Type:         "playing",
					OnlineStatus: "online",
				}
				if err := maintenance.SetBotStatus(s, statusData); err != nil {
					log.Printf("Error setting default bot status: %v", err)
					return
				}
				newSetting := &model.BotSystemSetting{ID: "status"}
				if err := newSetting.Value.Set(statusData); err != nil {
					log.Printf("Error setting default status JSON: %v", err)
					return
				}
				if err := settingsRepo.Create(newSetting); err != nil {
					log.Printf("Error creating default bot status: %v", err)
				}
				return
			}
			botStatus := botStatusRecord.Value

			// jsonからデータを引き出す処理を追加
			var statusData maintenance.StatusData
			err = json.Unmarshal(botStatus.Bytes, &statusData)

			if err != nil {
				log.Printf("Error unmarshaling bot status: %v", err)
				return
			} else {
				// Set Bot Status
				err = maintenance.SetBotStatus(s, statusData)
				if err != nil {
					log.Printf("Error setting bot status: %v", err)
					return
				}
			}
		}
	}
}
