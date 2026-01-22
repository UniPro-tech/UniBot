package handler

import (
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(ctx *internal.BotContext) func(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
		if s.VoiceConnections[vsu.GuildID] == nil {
			return
		}

		if vsu.Member.User.Bot {
			return
		}

		changeType := "left"
		if vsu.BeforeUpdate == nil || vsu.BeforeUpdate.ChannelID != s.VoiceConnections[vsu.GuildID].ChannelID {
			changeType = "joined"
		} else if vsu.BeforeUpdate != nil && vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != "" {
			changeType = "moved"
		}

		if changeType == "joined" && s.VoiceConnections[vsu.GuildID].ChannelID == vsu.ChannelID {
			channel, err := s.State.Channel(vsu.ChannelID)
			if err != nil {
				log.Printf("Error fetching channel: %v", err)
				return
			}
			text := fmt.Sprintf("%sが %s に参加しました。", vsu.Member.DisplayName(), channel.Name)
			voice.SynthesizeAndPlay(ctx, s, repository.DefaultTTSPersonalSetting, vsu.GuildID, text)
			return
		}
		if changeType == "left" && s.VoiceConnections[vsu.GuildID].ChannelID == vsu.BeforeUpdate.ChannelID {
			guild, err := s.State.Guild(vsu.GuildID)
			voiceStates := guild.VoiceStates
			var stillInChannel bool
			for _, vs := range voiceStates {
				if vs.UserID == vsu.Member.User.ID && vs.ChannelID == s.VoiceConnections[vsu.GuildID].ChannelID {
					stillInChannel = true
					break
				}
			}
			if !stillInChannel {
				s.VoiceConnections[vsu.GuildID].Disconnect()
				repo := repository.NewTTSConnectionRepository(ctx.DB)
				data, err := repo.GetByGuildID(vsu.GuildID)
				if err != nil {
					log.Printf("Error fetching TTS connection data: %v", err)
					return
				}
				if data != nil {
					channelId := data.ChannelID
					err = repo.DeleteByGuildID(vsu.GuildID)
					if err != nil {
						log.Printf("Error deleting TTS connection data: %v", err)
						return
					}
					s.ChannelMessageSendComplex(channelId,
						&discordgo.MessageSend{
							Embed: &discordgo.MessageEmbed{
								Title:       "TTS接続解除",
								Description: "ボイスチャンネルから誰もいなくなったため、TTSの接続を解除しました。",
								Color:       ctx.Config.Colors.Success,
								Timestamp:   time.Now().Format(time.RFC3339),
							},
						},
					)
				}
			}

			channel, err := s.State.Channel(vsu.BeforeUpdate.ChannelID)
			if err != nil {
				log.Printf("Error fetching channel: %v", err)
				return
			}
			text := fmt.Sprintf("%sが %s から退出しました。", vsu.Member.DisplayName(), channel.Name)
			voice.SynthesizeAndPlay(ctx, s, repository.DefaultTTSPersonalSetting, vsu.GuildID, text)
			return
		}
		if changeType == "moved" && s.VoiceConnections[vsu.GuildID].ChannelID == vsu.BeforeUpdate.ChannelID {
			channel, err := s.State.Channel(vsu.ChannelID)
			if err != nil {
				log.Printf("Error fetching channel: %v", err)
				return
			}
			text := fmt.Sprintf("%sが %s に移動しました。", vsu.Member.DisplayName(), channel.Name)
			voice.SynthesizeAndPlay(ctx, s, repository.DefaultTTSPersonalSetting, vsu.GuildID, text)
			return
		}
	}
}
