package maintenance

import (
	"context"
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/query"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
)

type StatusData struct {
	Text         string
	Type         discord.ActivityType
	OnlineStatus discord.OnlineStatus
}

func LoadStatusCommandContext() discord.ApplicationCommandOptionSubCommandGroup {
	return discord.ApplicationCommandOptionSubCommandGroup{
		Name:        "status",
		Description: "Botのステータスを更新します。",
		Options: []discord.ApplicationCommandOptionSubCommand{
			{
				Name:        "set",
				Description: "Botのステータスを設定します",
				Options: []discord.ApplicationCommandOption{
					discord.ApplicationCommandOptionString{
						Name:        "text",
						Required:    true,
						Description: "ステータスメッセージ",
					},
					discord.ApplicationCommandOptionInt{
						Name:        "status",
						Description: "ステータス種類",
						Required:    true,
						Choices: []discord.ApplicationCommandOptionChoiceInt{
							{Name: "playing", Value: int(discordgo.ActivityTypeGame)},
							{Name: "streaming", Value: int(discordgo.ActivityTypeStreaming)},
							{Name: "listening", Value: int(discordgo.ActivityTypeListening)},
							{Name: "watching", Value: int(discordgo.ActivityTypeWatching)},
							{Name: "competing", Value: int(discordgo.ActivityTypeCompeting)},
							{Name: "custom", Value: int(discordgo.ActivityTypeCustom)},
						},
					},
					discord.ApplicationCommandOptionString{
						Name:        "type",
						Description: "オンライン状態",
						Required:    true,
						Choices: []discord.ApplicationCommandOptionChoiceString{
							{Name: "online", Value: "online"},
							{Name: "idle", Value: "idle"},
							{Name: "dnd", Value: "dnd"},
							{Name: "invisible", Value: "invisible"},
						},
					},
				},
			}, {
				Name:        "reset",
				Description: "Botのステータスをリセットします",
			},
		},
	}
}

func StatusResetHandler(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		err := ResetBotStatus(e.Client())
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ステータスのリセットに失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}

		// DB Reset
		systemPrefData, err := query.SystemPreference.FirstOrInit()
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ステータス設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}
		systemPrefData.StatusType = "online"
		systemPrefData.ActivitySummary = nil
		systemPrefData.ActivityType = nil
		err = query.SystemPreference.Save(systemPrefData)
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ステータス設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}

		responseEmbed := discord.Embed{
			Title:       "ステータスリセット",
			Description: "Botのステータスをデフォルトにリセットしました。",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}

func StatusSetHandler(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		var statusText string
		var onlineStatus discord.OnlineStatus
		var statusType discord.ActivityType
		for _, option := range data.Options {
			switch option.Name {
			case "text":
				if option.Type == discord.ApplicationCommandOptionTypeString {
					statusText = string(option.Value)
				}
			case "status":
				if option.Type == discord.ApplicationCommandOptionTypeInt {
					statusType = discord.ActivityType(option.Int())
				}
			case "type":
				if option.Type == discord.ApplicationCommandOptionTypeString {
					onlineStatus = discord.OnlineStatus(option.Value)
				}
			}
		}

		err := SetBotStatus(e.Client(), StatusData{Text: statusText, Type: statusType, OnlineStatus: onlineStatus})
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ステータスのセットに失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}

		systemPrefData, err := query.SystemPreference.FirstOrInit()
		if err != nil {
			errorEmbed := discord.Embed{
				Title:       "エラー",
				Description: "設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed))
			return err
		}

		systemPrefData.ActivitySummary = &statusText
		statusTypeStr := ActivityTypeToString(statusType)
		systemPrefData.ActivityType = &statusTypeStr
		systemPrefData.StatusType = string(onlineStatus)

		err = query.SystemPreference.Save(systemPrefData)
		if err != nil {
			errorEmbed := discord.Embed{
				Title:       "エラー",
				Description: "設定の更新に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed))
			return err
		}

		responseEmbed := discord.Embed{
			Title:       "ステータス更新",
			Description: "Botのステータスを更新しました。",
			Color:       config.Colors.Success,
			Fields: []discord.EmbedField{
				{
					Name:  "ステータスメッセージ",
					Value: statusText,
				},
				{
					Name:  "ステータス種類",
					Value: statusTypeStr,
				},
				{
					Name:  "オンライン状態",
					Value: string(onlineStatus),
				},
			},
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}

func ActivityTypeToString(activityType discord.ActivityType) string {
	switch activityType {
	case discord.ActivityTypeGame:
		return "playing"
	case discord.ActivityTypeStreaming:
		return "streaming"
	case discord.ActivityTypeListening:
		return "listening"
	case discord.ActivityTypeWatching:
		return "watching"
	case discord.ActivityTypeCompeting:
		return "competing"
	case discord.ActivityTypeCustom:
		return "custom"
	default:
		return "unknown"
	}
}

func StringToActivityType(typeStr string) discord.ActivityType {
	switch typeStr {
	case "game":
		return discord.ActivityTypeGame
	case "streaming":
		return discord.ActivityTypeStreaming
	case "listening":
		return discord.ActivityTypeListening
	case "watching":
		return discord.ActivityTypeWatching
	case "competing":
		return discord.ActivityTypeCompeting
	case "custom":
		return discord.ActivityTypeCustom
	default:
		return discord.ActivityTypeGame
	}
}

func StringToOnlineStatus(typeStr string) discord.OnlineStatus {
	switch typeStr {
	case "online":
		return discord.OnlineStatusOnline
	case "dnd":
		return discord.OnlineStatusDND
	case "idle":
		return discord.OnlineStatusIdle
	case "invisible":
		return discord.OnlineStatusInvisible
	case "offline":
		return discord.OnlineStatusOffline
	default:
		return discord.OnlineStatusOnline
	}
}

func SetBotStatus(client *bot.Client, data StatusData) error {
	return client.SetPresence(context.Background(), gateway.PresenceOpt(func(p *gateway.MessageDataPresenceUpdate) {
		p.Activities = []discord.Activity{
			{
				Type: data.Type,
				Name: data.Text,
			},
		}
		p.Status = data.OnlineStatus
		p.AFK = false
	}))
}

func ResetBotStatus(client *bot.Client) error {
	systemPreference, err := query.SystemPreference.FirstOrInit()
	statusData := StatusData{
		Text: *systemPreference.ActivitySummary,
	}
	if err != nil {
		log.Printf("Error setting bot status: %v", err)
	} else {
		systemPreference.StatusType = "online"
		t := "game"
		systemPreference.ActivityType = &t
		serverCounts := client.Caches.GuildsLen()
		s := fmt.Sprintf("Serving %d servers | /help", serverCounts)
		systemPreference.ActivitySummary = &s
		err := query.SystemPreference.Save(systemPreference)
		if err != nil {
			log.Printf("Error setting bot status: %v", err)
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
	}
	err = SetBotStatus(client, statusData)
	if err != nil {
		log.Printf("Error setting bot status: %v", err)
	}
	return err
}
