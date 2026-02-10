package messageComponent

import (
	"unibot/internal"
	schedulecmd "unibot/internal/bot/command/general/schedule"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterHandler(schedulecmd.ScheduleModalOnetimeButtonID, HandleScheduleOnetimeModal)
	RegisterHandler(schedulecmd.ScheduleModalRepeatButtonID, HandleScheduleRepeatModal)
}

func HandleScheduleOnetimeModal(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	schedulecmd.ShowOnetimeModal(s, i)
}

func HandleScheduleRepeatModal(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	schedulecmd.ShowRepeatModal(s, i)
}
