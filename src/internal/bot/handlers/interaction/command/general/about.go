package general

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"

	"unibot/internal"
)

func LoadAboutCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "about",
		Description: "ボットの情報を表示します",
	}
}

func About(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		// コントリビューター一覧をMarkdown形式で作成
		contributorsText := ""
		// Botは最後にする
		for _, contributor := range config.Contributors {
			if contributor.IsBot {
				continue
			}
			contributorsText += "- [" + contributor.Username + "](" + contributor.Profile + ")\n"
		}
		for _, contributor := range config.Contributors {
			if !contributor.IsBot {
				continue
			}
			contributorsText += "- [" + contributor.Username + "](" + contributor.Profile + ")\n"
		}

		responseEmbed := discord.Embed{
			Title:       "About " + config.BotName + " 🤖",
			Description: config.Description,
			Color:       config.Colors.Primary,
			Fields: []discord.EmbedField{
				{
					Name:  "Version",
					Value: config.BotVersion,
				},
				{
					Name:  "Contributors",
					Value: contributorsText,
				},
				{
					Name:  "GitHub",
					Value: config.GitHub,
				},
				{
					Name:  "Support Server",
					Value: config.SupportServer,
				},
			},
			Footer: &discord.EmbedFooter{
				Text:    "Requested by " + e.User().Username,
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}
