package dict

import (
	"fmt"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadListCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "list",
		Description: "TTS辞書の単語一覧を表示します",
	}
}

func List(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		repo := repository.NewTTSDictionaryRepository(ctx.DB)

		entries, err := repo.ListByGuild(e.GuildID().String())
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "辞書の取得中にエラーが発生しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		if len(entries) == 0 {
			responseEmbed := discord.Embed{
				Title:       "辞書が空です",
				Description: "辞書に登録されている単語がありません。",
				Color:       config.Colors.Warning,
				Footer: &discord.EmbedFooter{
					Text:    fmt.Sprintf("Requested by %s", e.User().Username),
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		// 辞書エントリをフォーマット
		var lines []string
		for _, entry := range entries {
			caseMark := ""
			if entry.CaseSensitive {
				caseMark = " [大小区別]"
			}
			lines = append(lines, fmt.Sprintf("• **%s** → %s%s", entry.Word, entry.Definition, caseMark))
		}

		description := strings.Join(lines, "\n")
		if len(description) > 4000 {
			description = description[:4000] + "\n..."
		}

		responseEmbed := discord.Embed{
			Title:       fmt.Sprintf("TTS辞書 (%d件)", len(entries)),
			Description: description,
			Color:       config.Colors.Primary,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
		return err
	}
}
