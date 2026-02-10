package schedule

import (
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

const (
	ScheduleModalOnetimeButtonID = "schedule_open_onetime"
	ScheduleModalRepeatButtonID  = "schedule_open_repeat"
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
	hasRepeatOption := false

	for _, opt := range options {
		if opt.Name == "repeat" {
			isRepeat = opt.BoolValue()
			hasRepeatOption = true
		}
	}

	showOnetime := true
	showRepeat := false
	if hasRepeatOption {
		showOnetime = !isRepeat
		showRepeat = isRepeat
	}

	promptScheduleModal(s, i, showOnetime, showRepeat)
}

func promptScheduleModal(s *discordgo.Session, i *discordgo.InteractionCreate, showOnetime, showRepeat bool) {
	components := buildScheduleModalButtons(showOnetime, showRepeat)
	if len(components) == 0 {
		components = buildScheduleModalButtons(true, false)
	}

	content := "単発の予約投稿を作成します。下のボタンからフォームを開いてください。"
	if showOnetime && showRepeat {
		content = "作成する予約投稿の種類を選んでください。"
	} else if showRepeat {
		content = "繰り返しの予約投稿を作成します。下のボタンからフォームを開いてください。"
	}

	_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
		Content:    content,
		Components: components,
		Flags:      discordgo.MessageFlagsEphemeral,
	})
}

func buildScheduleModalButtons(showOnetime, showRepeat bool) []discordgo.MessageComponent {
	buttons := make([]discordgo.MessageComponent, 0, 2)
	if showOnetime {
		buttons = append(buttons, discordgo.Button{
			CustomID: ScheduleModalOnetimeButtonID,
			Label:    "単発",
			Style:    discordgo.PrimaryButton,
		})
	}
	if showRepeat {
		buttons = append(buttons, discordgo.Button{
			CustomID: ScheduleModalRepeatButtonID,
			Label:    "繰り返し",
			Style:    discordgo.SecondaryButton,
		})
	}
	if len(buttons) == 0 {
		return nil
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{Components: buttons},
	}
}

// ShowOnetimeModal opens the one-time schedule modal.
func ShowOnetimeModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

// ShowRepeatModal opens the repeating schedule modal.
func ShowRepeatModal(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "schedule_create_repeat",
			Title:    "予約投稿の作成",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "time",
						Label:       "投稿時間 (JST)",
						Style:       discordgo.TextInputShort,
						Placeholder: "例: every day at 12:00, every Monday at 09:00",
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
