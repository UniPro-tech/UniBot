package handler

import (
	"log"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(ctx *internal.BotContext) func(s *discordgo.Session, r *discordgo.MessageCreate) {
	return func(s *discordgo.Session, r *discordgo.MessageCreate) {

		// Ignore bot itself
		if r.Author.ID == s.State.User.ID {
			return
		}

		// Ignore DM
		if r.GuildID == "" {
			return
		}

		// ----- TTS -----

		repo := repository.NewTTSConnectionRepository(ctx.DB)

		ttsConnectionData, err := repo.GetByGuildID(r.GuildID)
		if err != nil {
			log.Println(err)
			return
		}

		if ttsConnectionData != nil {
			userID := r.Author.ID

			if r.Author.Bot {
				return
			}

			if r.ChannelID != ttsConnectionData.ChannelID && r.ChannelID != s.VoiceConnections[r.GuildID].ChannelID {
				return
			}

			personalSetting, err := repository.NewTTSPersonalSettingRepository(ctx.DB).GetByMember(userID)
			if err != nil {
				log.Println(err)
				return
			}
			if personalSetting == nil {
				personalSetting = &repository.DefaultTTSPersonalSetting
			}

			go voice.SynthesizeAndPlay(ctx, s, *personalSetting, r.GuildID, r.Content)
		}
	}
}
