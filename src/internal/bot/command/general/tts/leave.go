package tts

import (
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadLeaveCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "leave",
		Description: "ボイスチャンネルから退出します",
	}
}

func Leave(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	userVoiceState, err := s.State.VoiceState(i.GuildID, i.Member.User.ID)
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルの情報を取得できませんでした。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			})
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)
	if userVoiceState == nil || userVoiceState.ChannelID == "" {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "先にボイスチャンネルに参加してください。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	botVoiceStatus, err := s.State.VoiceState(i.GuildID, s.State.User.ID)
	if err != nil && err != discordgo.ErrStateNotFound {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "Botの情報を取得できませんでした。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}
	if botVoiceStatus == nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルに参加していません。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}
	if botVoiceStatus.ChannelID != userVoiceState.ChannelID {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "同じボイスチャンネルに参加していません。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	voiceConnection := s.VoiceConnections[i.GuildID]
	if voiceConnection == nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルに接続していません。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	err = voiceConnection.Disconnect()
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルから退出できませんでした。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	dbConnection := ctx.DB
	repo := repository.NewTTSConnectionRepository(dbConnection)

	mgr := voice.GetManager()
	player := mgr.Get(i.GuildID)

	if player != nil {
		player.Close()
		if vc := player.GetVC(); vc != nil {
			vc.Disconnect()
		}
		mgr.Delete(i.GuildID)
	}

	err = repo.DeleteByGuildID(i.GuildID)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "TTS接続情報の削除に失敗しました。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "TTSボイスチャンネル退出",
					Description: "ボイスチャンネルから退出しました。",
					Color:       config.Colors.Success,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
}
