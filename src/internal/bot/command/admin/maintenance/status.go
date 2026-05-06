package maintenance

import (
	"context"
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

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
		Name: "status",
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
					IconURL: *e.Member().Avatar,
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
		database := ctx.DB
		repo := repository.NewBotSystemSettingRepository(database)
		err = repo.Delete("status")
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "ステータス設定の削除に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: *e.Member().Avatar,
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
				IconURL: *e.Member().Avatar,
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
				Description: "ステータスのリセットに失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: *e.Member().Avatar,
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}

		statusTypeStr := activityTypeToString(statusType)

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
				IconURL: *e.Member().Avatar,
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		database := ctx.DB
		repo := repository.NewBotSystemSettingRepository(database)
		listSettings, err := repo.List()
		if err != nil {
			errorEmbed := discord.Embed{
				Title:       "エラー",
				Description: "設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: *e.Member().Avatar,
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed))
			return err
		}

		// list settings の中に status というキーがあれば更新、なければ新規作成
		var statusSetting *model.BotSystemSetting
		for _, setting := range listSettings {
			if setting.ID == "status" {
				statusSetting = setting
				break
			}
		}
		if statusSetting != nil {
			statusSetting.Value.Set(StatusData{Text: statusText, Type: statusType, OnlineStatus: onlineStatus})
			err = repo.Update(statusSetting)
			if err != nil {
				errorEmbed := discord.Embed{
					Title:       "エラー",
					Description: "設定の更新に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", e.User().Username),
						IconURL: *e.Member().Avatar,
					},
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}
				_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed))
				return err
			}
		} else {
			newSetting := &model.BotSystemSetting{
				ID: "status",
			}
			newSetting.Value.Set(StatusData{Text: statusText, Type: statusType, OnlineStatus: onlineStatus})
			err = repo.Create(newSetting)
			if err != nil {
				errorEmbed := discord.Embed{
					Title:       "エラー",
					Description: "設定の更新に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						Text:    fmt.Sprintf("Requested by %s", e.User().Username),
						IconURL: *e.Member().Avatar,
					},
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}
				_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed))
				return err
			}
		}

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}

func activityTypeToString(activityType discord.ActivityType) string {
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

func stringToActivityType(typeStr string) discordgo.ActivityType {
	switch typeStr {
	case "playing":
		return discordgo.ActivityTypeGame
	case "streaming":
		return discordgo.ActivityTypeStreaming
	case "listening":
		return discordgo.ActivityTypeListening
	case "watching":
		return discordgo.ActivityTypeWatching
	case "competing":
		return discordgo.ActivityTypeCompeting
	case "custom":
		return discordgo.ActivityTypeCustom
	default:
		return discordgo.ActivityTypeGame
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
	serverCounts := client.Caches.GuildCache().Len()
	return client.SetPresence(context.Background(), gateway.PresenceOpt(func(p *gateway.MessageDataPresenceUpdate) {
		p.Activities = []discord.Activity{
			{
				Type: discord.ActivityTypeGame,
				Name: fmt.Sprintf("Serving %d servers | /help", serverCounts),
			},
		}
		p.Status = discord.OnlineStatusOnline
		p.AFK = false
	}))
}
