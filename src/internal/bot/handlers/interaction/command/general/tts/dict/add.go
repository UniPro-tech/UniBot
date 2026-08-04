package dict

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/query"
	"unibot/internal/util"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

func LoadAddCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "add",
		Description: "TTS辞書に単語を追加します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "word",
				Description: "追加する単語",
				Required:    true,
			},
			discord.ApplicationCommandOptionString{
				Name:        "yomi",
				Description: "追加する単語の読み",
				Required:    true,
			},
		},
	}
}

func Add(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		options := data.Options

		var word, yomi string

		for _, opt := range options {
			switch opt.Name {
			case "word":
				word = opt.String()
			case "yomi":
				yomi = opt.String()
			}
		}

		// 既存のエントリがあるか確認
		existing, err := query.TtsDictionary.Where(query.TtsDictionary.Word.Eq(word), query.TtsDictionary.GuildID.Eq(int64(*e.GuildID()))).Find()
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "辞書の確認中にエラーが発生しました。",
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
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		if existing != nil || len(existing) > 0 {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "この単語はすでに辞書に存在します。",
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
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed).WithEphemeral(true))
			return err
		}

		// 新しいエントリを作成
		entry := &model.TtsDictionary{
			GuildID: int64(*e.GuildID()),
			UserID:  int64(e.User().ID),
			Word:    word,
			Yomi:    yomi,
		}

		err = query.TtsDictionary.Create(entry)
		if err != nil {
			responseEmbed := discord.Embed{
				Title:       "エラー",
				Description: "辞書への追加中にエラーが発生しました。",
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
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
			return err
		}

		// 辞書キャッシュを無効化
		util.GetDictionaryCache().Invalidate(e.GuildID().String())

		responseEmbed := discord.Embed{
			Title: "単語を辞書に追加しました！",
			Color: config.Colors.Success,
			Fields: []discord.EmbedField{
				{
					Name:  "単語",
					Value: word,
					Inline: func() *bool {
						v := true
						return &v
					}(),
				},
				{
					Name:  "読み",
					Value: yomi,
					Inline: func() *bool {
						v := true
						return &v
					}(),
				},
			},
			Footer: &discord.EmbedFooter{
				Text:    fmt.Sprintf("Requested by %s", e.User().Username),
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(responseEmbed))
		return err
	}
}
