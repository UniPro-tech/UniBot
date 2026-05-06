package dict

import (
	"fmt"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"
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
				Name:        "definition",
				Description: "追加する単語の読み",
				Required:    true,
			},
			discord.ApplicationCommandOptionBool{
				Name:        "case_sensitive",
				Description: "大文字小文字を区別するか (デフォルト: false)",
				Required:    false,
			},
		},
	}
}

func Add(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		options := data.Options

		var word, definition string
		caseSensitive := false

		for _, opt := range options {
			switch opt.Name {
			case "word":
				word = opt.String()
			case "definition":
				definition = opt.String()
			case "case_sensitive":
				caseSensitive = opt.Bool()
			}
		}

		repo := repository.NewTTSDictionaryRepository(ctx.DB)

		// 既存のエントリがあるか確認
		existing, err := repo.GetByGuildWord(e.GuildID().String(), word)
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

		if existing != nil {
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
		entry := &model.TTSDictionary{
			GuildID:       e.GuildID().String(),
			UserID:        e.User().ID.String(),
			Word:          word,
			Definition:    definition,
			CaseSensitive: caseSensitive,
		}

		err = repo.Create(entry)
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
					Value: definition,
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
