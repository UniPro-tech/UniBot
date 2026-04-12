package handler

import (
	"context"
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

		if vsu.BeforeUpdate != nil && vsu.ChannelID != "" && vsu.ChannelID == vsu.BeforeUpdate.ChannelID {
			return
		}

		changeType := "left"
		botChannelID := getBotChannelID(s, vsu.GuildID)

		if vsu.BeforeUpdate == nil || vsu.BeforeUpdate.ChannelID != botChannelID {
			changeType = "joined"
		} else if vsu.BeforeUpdate != nil && vsu.ChannelID != "" && vsu.BeforeUpdate.ChannelID != "" {
			changeType = "moved"
		}

		if changeType == "joined" && botChannelID == vsu.ChannelID {
			channel, err := s.State.Channel(vsu.ChannelID)
			if err != nil {
				log.Printf("Error fetching channel: %v", err)
				return
			}
			text := fmt.Sprintf("%sが %s に参加しました。", vsu.Member.DisplayName(), channel.Name)
			vp := voice.GetManager().GetOrCreate(vsu.GuildID, botChannelID, s.VoiceConnections[vsu.GuildID], ctx)
			vp.EnqueueText(voice.QueueItem{
				Text:    text,
				Setting: repository.DefaultTTSPersonalSetting,
			})
			return
		}
		if changeType == "left" && botChannelID == vsu.BeforeUpdate.ChannelID {
			guild, err := s.State.Guild(vsu.GuildID)
			voiceStates := guild.VoiceStates
			var stillInChannel bool
			for _, vs := range voiceStates {
				user, err := s.User(vs.UserID)
				if err != nil || user.Bot {
					continue
				}
				if vs.ChannelID == botChannelID {
					stillInChannel = true
					break
				}
			}
			if !stillInChannel {
				backCtx := context.Background()
				s.VoiceConnections[vsu.GuildID].Disconnect(backCtx)
				repo := repository.NewTTSConnectionRepository(ctx.DB)
				data, err := repo.GetByGuildID(vsu.GuildID)
				if err != nil {
					log.Printf("Error fetching TTS connection data: %v", err)
					return
				}
				mgr := voice.GetManager()
				player := mgr.Get(vsu.GuildID)

				if player != nil {
					player.Close()
					if vc := player.GetVC(); vc != nil {
						vc.Disconnect(backCtx)
					}
					mgr.Delete(vsu.GuildID)
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
			vp := voice.GetManager().GetOrCreate(
				vsu.GuildID,
				botChannelID,
				s.VoiceConnections[vsu.GuildID],
				ctx,
			)
			vp.EnqueueText(voice.QueueItem{
				Text:    text,
				Setting: repository.DefaultTTSPersonalSetting,
			})
			return
		}
		if changeType == "moved" && botChannelID == vsu.BeforeUpdate.ChannelID {
			channel, err := s.State.Channel(vsu.ChannelID)
			if err != nil {
				log.Printf("Error fetching channel: %v", err)
				return
			}
			text := fmt.Sprintf("%sが %s に移動しました。", vsu.Member.DisplayName(), channel.Name)
			vp := voice.GetManager().GetOrCreate(
				vsu.GuildID,
				botChannelID,
				s.VoiceConnections[vsu.GuildID],
				ctx,
			)
			vp.EnqueueText(voice.QueueItem{
				Text:    text,
				Setting: repository.DefaultTTSPersonalSetting,
			})
			return
		}
	}
}

func getBotChannelID(s *discordgo.Session, guildID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		return ""
	}

	for _, vs := range guild.VoiceStates {
		if vs.UserID == s.State.User.ID {
			return vs.ChannelID
		}
	}
	return ""
}
