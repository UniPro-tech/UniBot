package maintenance

import (
	"fmt"
	"log"
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

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

		s.UpdateStatusComplex(discordgo.UpdateStatusData{
			Activities: []*discordgo.Activity{
				{
					Name: statusText,
					Type: statusType,
				},
			},
			Status: onlineStatus,
		})

		var statusTypeStr string
		switch statusType {
		case discordgo.ActivityTypeGame:
			statusTypeStr = "playing"
		case discordgo.ActivityTypeStreaming:
			statusTypeStr = "streaming"
		case discordgo.ActivityTypeListening:
			statusTypeStr = "listening"
		case discordgo.ActivityTypeWatching:
			statusTypeStr = "watching"
		case discordgo.ActivityTypeCompeting:
			statusTypeStr = "competing"
		case discordgo.ActivityTypeCustom:
			statusTypeStr = "custom"
		default:
			statusTypeStr = "unknown"
		}

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

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "ステータスを更新しました。",
				Embeds:  []*discordgo.MessageEmbed{responseEmbed},
			},
		})
	case "reset":
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
		err := s.UpdateStatusComplex(*defaultStatus)
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
