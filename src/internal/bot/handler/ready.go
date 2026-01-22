package handler

import (
	"encoding/json"
	"log"

	"unibot/internal"
	"unibot/internal/bot/command/admin/maintenance"
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
				maintenance.ResetBotStatus(s)
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
