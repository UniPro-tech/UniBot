package messageComponent

import (
	"fmt"
	"strconv"
	"time"
	"unibot/internal"
	"unibot/internal/repository"
	"unibot/internal/util"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// HandleTTSDictRemove は辞書削除のセレクトメニューを処理します
func HandleTTSDictRemove(ctx *internal.BotContext) func(data discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
	return func(_ discord.SelectMenuInteractionData, e *handler.ComponentEvent) error {
		config := ctx.Config
		data := e.StringSelectMenuInteractionData()
		selected := data.Values[0]

		selectedID, err := strconv.ParseUint(selected, 10, 64)

		repo := repository.NewTTSDictionaryRepository(ctx.DB)

		// 削除対象の単語を取得
		entry, err := repo.GetByID(uint(selectedID))
		if err != nil || entry == nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "単語が見つかりませんでした。",
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

		// セキュリティチェック: エントリが現在のギルドに属しているか確認
		if entry.GuildID != e.GuildID().String() {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "この単語を削除する権限がありません。",
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

		word := entry.Word

		// 削除実行
		err = repo.DeleteByID(uint(selectedID))
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "単語の削除に失敗しました。",
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

		// 辞書キャッシュを無効化
		util.GetDictionaryCache().Invalidate(entry.GuildID)

		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(discord.Embed{
			Title:       "単語を削除しました",
			Description: "「" + word + "」を辞書から削除しました。",
			Color:       config.Colors.Success,
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}).WithEphemeral(true))
		return err
	}
}
