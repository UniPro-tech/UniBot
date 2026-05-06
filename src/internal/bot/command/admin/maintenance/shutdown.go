package maintenance

import (
	"fmt"
	"os"
	"syscall"
	"time"
	"unibot/internal"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadShutdownCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "shutdown",
		Description: "Botプロセスをシャットダウンします。(コンテナの場合は再起動を兼ねます)",
	}
}

func Shutdown(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		responseEmbed := discord.Embed{
			Title:       "Now Shutting down...",
			Description: "The bot is Shutting down...",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))

		// 再起動と称してプロセスを終了する
		// コンテナで起動することを想定しているため、再起動となる
		go func() {
			time.Sleep(2 * time.Second)
			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				return
			}
			p.Signal(syscall.SIGTERM)
		}()

		return err
	}
}
