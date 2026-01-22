package handler

import (
	"fmt"
	"log"
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
