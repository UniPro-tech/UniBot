package tts

import (
	"unibot/internal"
	"unibot/internal/bot/voice"

	"github.com/bwmarrin/discordgo"
)

func LoadSkipCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "skip",
		Description: "現在再生中の音声をスキップします",
	}
}

func Skip(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	player := voice.GetManager().Get(i.GuildID)
	if player == nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "現在再生中の音声はありません。",
						Color:       ctx.Config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	player.SkipCurrent()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "音声スキップ",
					Description: "現在再生中の音声をスキップしました。",
					Color:       ctx.Config.Colors.Success,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		},
	})
}
