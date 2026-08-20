package event_handlers

import (
	"context"
	"fmt"
	"log/slog"

	"unibot/internal"
	"unibot/internal/bot/handlers/interaction/command/admin/maintenance"
	"unibot/internal/bot/notify"
	"unibot/internal/logger"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/events"
)

func Ready(ctx context.Context, bctx *internal.BotContext, e *events.Ready) {
	guildCount := len(e.Guilds)

	logger.Notice(ctx, "bot ready",
		slog.String("event", "ready"),
		slog.String("user", fmt.Sprintf("%s#%s", e.User.Username, e.User.Discriminator)),
		slog.Int("guilds", guildCount),
		slog.String("version", bctx.Config.BotVersion),
	)

	notify.Ready(ctx, e.Client(), bctx.Config, notify.ReadyInfo{
		User:   e.User.Username,
		Guilds: guildCount,
	})

	systemPreference, err := query.SystemPreference.FirstOrCreate()
	if err != nil {
		slog.ErrorContext(ctx, "failed to load system preference", slog.Any("err", err))
		return
	}

	if systemPreference.ActivityType == nil {
		t := "game"
		systemPreference.ActivityType = &t
	}
	if systemPreference.ActivitySummary == nil {
		s := fmt.Sprintf("Serving %d servers | /help", guildCount)
		systemPreference.ActivitySummary = &s
	}

	statusData := maintenance.StatusData{
		Text: *systemPreference.ActivitySummary,
	}

	if err := query.SystemPreference.Save(systemPreference); err != nil {
		slog.ErrorContext(ctx, "failed to save system preference", slog.Any("err", err))
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
		slog.ErrorContext(ctx, "failed to set bot presence", slog.Any("err", err))
		return
	}
}
