package maintenance

import (
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/db"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

type StatusData struct {
	Text         string `json:"text"`
	Type         string `json:"type"`
	OnlineStatus string `json:"online_status"`
}

func Status(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()

	options := i.ApplicationCommandData().Options[0].Options
	if len(options) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "サブコマンドが指定されていません。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	subCommand := options[0]
	switch subCommand.Name {
	case "set":
		var statusText, onlineStatus string
		var statusType discordgo.ActivityType
		for _, option := range subCommand.Options {
			switch option.Name {
			case "text":
				if option.Type == discordgo.ApplicationCommandOptionString {
					statusText = option.StringValue()
				}
			case "status":
				if option.Type == discordgo.ApplicationCommandOptionInteger {
					statusType = discordgo.ActivityType(option.IntValue())
				}
			case "type":
				if option.Type == discordgo.ApplicationCommandOptionString {
					onlineStatus = option.StringValue()
				}
			}
		}

		err := SetBotStatus(s, StatusData{Text: statusText, Type: activityTypeToString(statusType), OnlineStatus: onlineStatus})
		if err != nil {
			log.Fatalf("Failed to set status: %v", err)
			embed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "ステータスの設定に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		statusTypeStr := activityTypeToString(statusType)

		responseEmbed := &discordgo.MessageEmbed{
			Title:       "ステータス更新",
			Description: "Botのステータスを更新しました。",
			Color:       config.Colors.Success,
			Fields: []*discordgo.MessageEmbedField{
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
					Value: onlineStatus,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Requested by " + i.Member.DisplayName(),
				IconURL: i.Member.AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}

		database, err := db.NewDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			errorEmbed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "データベースへの接続に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{errorEmbed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		repo := repository.NewBotSystemSettingRepository(database)
		listSettings, err := repo.List()
		if err != nil {
			log.Printf("Failed to list settings: %v", err)
			errorEmbed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "設定の取得に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{errorEmbed},
					Flags:  discordgo.MessageFlagsEphemeral,
				},
			})
			return
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
			statusSetting.Value.Set(StatusData{Text: statusText, Type: statusTypeStr, OnlineStatus: onlineStatus})
			err = repo.Update(statusSetting)
			if err != nil {
				log.Printf("Failed to update status setting: %v", err)
				errorEmbed := &discordgo.MessageEmbed{
					Title:       "エラー",
					Description: "ステータス設定の更新に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{errorEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
		} else {
			newSetting := &model.BotSystemSetting{
				ID: "status",
			}
			newSetting.Value.Set(StatusData{Text: statusText, Type: statusTypeStr, OnlineStatus: onlineStatus})
			err = repo.Create(newSetting)
			if err != nil {
				log.Printf("Failed to create status setting: %v", err)
				errorEmbed := &discordgo.MessageEmbed{
					Title:       "エラー",
					Description: "ステータス設定の作成に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{errorEmbed},
						Flags:  discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
		}

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ステータスを更新しました。",
				Embeds:  []*discordgo.MessageEmbed{responseEmbed},
			},
		})
	case "reset":
		err := ResetBotStatus(s)
		if err != nil {
			log.Fatalf("Failed to reset status: %v", err)
			embed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "ステータスのリセットに失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		// DB Reset
		database, err := db.NewDB()
		if err != nil {
			log.Printf("Failed to connect to database: %v", err)
			embed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "データベースへの接続に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}
		repo := repository.NewBotSystemSettingRepository(database)
		err = repo.Delete("status")
		if err != nil {
			log.Printf("Failed to delete status setting: %v", err)
			embed := &discordgo.MessageEmbed{
				Title:       "エラー",
				Description: "ステータス設定の削除に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discordgo.MessageEmbedFooter{
					Text:    "Requested by " + i.Member.DisplayName(),
					IconURL: i.Member.AvatarURL(""),
				},
				Timestamp: time.Now().Format(time.RFC3339),
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{embed},
				},
			})
			return
		}

		embed := &discordgo.MessageEmbed{
			Title:       "ステータスリセット",
			Description: "Botのステータスをデフォルトにリセットしました。",
			Color:       config.Colors.Success,
			Footer: &discordgo.MessageEmbedFooter{
				Text:    "Requested by " + i.Member.DisplayName(),
				IconURL: i.Member.AvatarURL(""),
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
	default:
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "不明なサブコマンドです。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
	}
}

func activityTypeToString(activityType discordgo.ActivityType) string {
	switch activityType {
	case discordgo.ActivityTypeGame:
		return "playing"
	case discordgo.ActivityTypeStreaming:
		return "streaming"
	case discordgo.ActivityTypeListening:
		return "listening"
	case discordgo.ActivityTypeWatching:
		return "watching"
	case discordgo.ActivityTypeCompeting:
		return "competing"
	case discordgo.ActivityTypeCustom:
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

func SetBotStatus(s *discordgo.Session, data StatusData) error {
	log.Print("Updating Bot Status", data)
	return s.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: data.Text,
				Type: stringToActivityType(data.Type),
			},
		},
		Status: data.OnlineStatus,
	})
}

func ResetBotStatus(s *discordgo.Session) error {
	serverCounts := s.State.Guilds
	defaultStatus := &discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "Serving " + fmt.Sprintf("%d", len(serverCounts)) + " servers | /help",
				Type: discordgo.ActivityTypeGame,
			},
		},
		Status: "online",
	}
	return s.UpdateStatusComplex(*defaultStatus)
}
