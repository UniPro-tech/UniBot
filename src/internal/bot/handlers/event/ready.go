package event_handlers

import (
	"fmt"
	"log"

	"unibot/internal"
	"unibot/internal/bot/handlers/interaction/command/admin/maintenance"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func Ready(ctx *internal.BotContext, e *events.Ready) {
	log.Println("Bot is ready 🚀")
	log.Printf("Logged in as: %v#%v", e.User.Username, e.User.Discriminator)

	// Get status
	systemPreference, err := query.SystemPreference.FirstOrInit()
	if err != nil {
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

		// activity type switch
		activityType := "game"
		if systemPreference.ActivityType != nil {
			activityType = *systemPreference.ActivityType
		}
		switch activityType {
		case "custom":
			statusData.Type = discord.ActivityTypeCustom
		case "game":
			statusData.Type = discord.ActivityTypeGame
		case "competing":
			statusData.Type = discord.ActivityTypeCompeting
		case "watching":
			statusData.Type = discord.ActivityTypeWatching
		case "listening":
			statusData.Type = discord.ActivityTypeListening
		case "streaming":
			statusData.Type = discord.ActivityTypeStreaming
		default:
			statusData.Type = discord.ActivityTypeGame
		}

		// status type switch
		switch systemPreference.StatusType {
		case "online":
			statusData.OnlineStatus = discord.OnlineStatusOnline
		case "dnd":
			statusData.OnlineStatus = discord.OnlineStatusDND
		case "idle":
			statusData.OnlineStatus = discord.OnlineStatusIdle
		case "invisible":
			statusData.OnlineStatus = discord.OnlineStatusInvisible
		case "offline":
			statusData.OnlineStatus = discord.OnlineStatusOffline
		default:
			statusData.OnlineStatus = discord.OnlineStatusOnline
		}
		if err := maintenance.SetBotStatus(e.Client(), statusData); err != nil {
			log.Printf("Error setting bot status: %v", err)
			return
		}
	}
}
