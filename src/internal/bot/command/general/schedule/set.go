package schedule

import (
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

func LoadSetCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "set",
		Description: "予約投稿を作成します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "repeat",
				Description: "繰り返し投稿にするかどうか",
				Required:    false,
			},
		},
	}
}

func Set(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	isRepeat := false

	for _, opt := range options {
		if opt.Name == "repeat" {
			isRepeat = opt.BoolValue()
		}
	}

	if isRepeat {
		showRepeatModal(s, i)
		return
	}

	showOnetimeModal(s, i)
}

func showOnetimeModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "schedule_create_onetime",
			Title:    "予約投稿の作成",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "time",
						Label:       "投稿時間 (YYYY-MM-DD HH:mm / JST)",
						Style:       discordgo.TextInputShort,
						Placeholder: "例: 2026-12-31 23:59",
						Required:    true,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "message",
						Label:       "投稿内容",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "投稿内容を入力してください",
						Required:    true,
					},
				}},
			},
		},
	})
}

func showRepeatModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "schedule_create_repeat",
			Title:    "予約投稿の作成",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "time",
						Label:       "投稿時間 (cron / JST)",
						Style:       discordgo.TextInputShort,
						Placeholder: "例: 0 9 * * *",
						Required:    true,
					},
				}},
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "message",
						Label:       "投稿内容",
						Style:       discordgo.TextInputParagraph,
						Placeholder: "投稿内容を入力してください",
						Required:    true,
					},
				}},
			},
		},
	})
}
