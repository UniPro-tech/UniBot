package maintenance

import (
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

func Reboot(s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := internal.LoadConfig()

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Now Rebooting",
					Description: "The bot is rebooting...",
					Color:       config.Colors.Success,
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
