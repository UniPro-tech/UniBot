package dict

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadRemoveCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "remove",
		Description: "TTS辞書から単語を削除します",
	}
}

func Remove(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config
		repo := repository.NewTTSDictionaryRepository(ctx.DB)

		// 管理者かどうか確認
		perms := e.Member().Permissions
		isAdmin := perms&discordgo.PermissionAdministrator != 0

		// 辞書エントリを取得
		var entries []*model.TTSDictionary
		var err error
		if isAdmin {
			entries, err = repo.ListByGuild(e.GuildID().String())
		} else {
			entries, err = repo.ListByGuildUser(e.GuildID().String(), e.User().ID.String())
		}

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

		// 25件以上ある場合は最初の25件のみ表示
		displayEntries := entries
		if len(entries) > 25 {
			displayEntries = entries[:25]
		}

		// セレクトメニューを作成
		options := make([]discord.StringSelectMenuOption, len(displayEntries))
		for idx, entry := range displayEntries {
			options[idx] = discord.StringSelectMenuOption{
				Label:       entry.Word,
				Value:       fmt.Sprintf("%d", entry.ID),
				Description: entry.Definition,
			}
		}

		content := "削除したい単語を選んでください。"
		components := []discord.LayoutComponent{
			discord.ActionRowComponent{
				Components: []discord.InteractiveComponent{
					discord.StringSelectMenuComponent{
						CustomID:    "tts_dict_remove",
						Placeholder: "削除する単語を選んでください",
						Options:     options,
					},
				},
			},
		}

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithContent(content).WithComponents(components...))
		return err
	}
}
