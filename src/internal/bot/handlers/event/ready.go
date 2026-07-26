package event_handlers

import (
	"fmt"
	"log"

	"unibot/internal"
	"unibot/internal/bot/handlers/interaction/command/admin/maintenance"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/events"
)

func Ready(ctx *internal.BotContext, e *events.Ready) {
	log.Println("Bot is ready 🚀")
	log.Printf("Logged in as: %v#%v", e.User.Username, e.User.Discriminator)

	// Get status
	systemPreference, err := query.SystemPreference.FirstOrCreate()
	if err != nil {
		log.Printf("Error setting bot status: %v", err)
		return
	} else {
		if systemPreference.ActivityType == nil {
			t := "game"
			systemPreference.ActivityType = &t
		}
		if systemPreference.ActivitySummary == nil {
			serverCounts := len(e.Guilds)
			s := fmt.Sprintf("Serving %d servers | /help", serverCounts)
			systemPreference.ActivitySummary = &s
		}
		statusData := maintenance.StatusData{
			Text: *systemPreference.ActivitySummary,
		}
		if err := query.SystemPreference.Save(systemPreference); err != nil {
			log.Printf("Error setting bot status: %v", err)
			return
		}

		// activity type switch
		activityType := "game"
		if systemPreference.ActivityType != nil {
			activityType = *systemPreference.ActivityType
		}
		statusData.Type = maintenance.StringToActivityType(activityType)

		// status type switch
		statusData.OnlineStatus = maintenance.StringToOnlineStatus(systemPreference.StatusType)
		if err := maintenance.SetBotStatus(e.Client(), statusData); err != nil {
			log.Printf("Error setting bot status: %v", err)
			return
		}
	}
}
