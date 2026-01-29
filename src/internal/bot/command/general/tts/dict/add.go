package dict

import (
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"
	"unibot/internal/util"

	"github.com/bwmarrin/discordgo"
)

func LoadAddCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "add",
		Description: "TTS辞書に単語を追加します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "word",
				Description: "追加する単語",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "definition",
				Description: "追加する単語の読み",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "case_sensitive",
				Description: "大文字小文字を区別するか (デフォルト: false)",
				Required:    false,
			},
		},
	}
}

func Add(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ボイスチャンネルの情報を取得できませんでした。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			})
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)

	options := i.ApplicationCommandData().Options[0].Options[0].Options

	var word, definition string
	var caseSensitive bool

	for _, opt := range options {
		switch opt.Name {
		case "word":
			word = opt.StringValue()
		case "definition":
			definition = opt.StringValue()
		case "case_sensitive":
			caseSensitive = opt.BoolValue()
		}
	}

	repo := repository.NewTTSDictionaryRepository(ctx.DB)

	// 既存のエントリがあるか確認
	existing, err := repo.GetByGuildWord(i.GuildID, word)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "エラー",
						Description: "辞書の確認中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	if existing != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "エラー",
						Description: "この単語はすでに辞書に存在します。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// 新しいエントリを作成
	entry := &model.TTSDictionary{
		GuildID:       i.GuildID,
		UserID:        i.Member.User.ID,
		Word:          word,
		Definition:    definition,
		CaseSensitive: caseSensitive,
	}

	err = repo.Create(entry)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
						Title:       "エラー",
						Description: "辞書への追加中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	// 辞書キャッシュを無効化
	util.GetDictionaryCache().Invalidate(i.GuildID)

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{
			{
					Title: "単語を辞書に追加しました！",
					Color: config.Colors.Success,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "単語",
							Value:  word,
							Inline: true,
						},
						{
							Name:   "読み",
							Value:  definition,
							Inline: true,
						},
					},
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
			},
		},
	})
		if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "辞書への追加後の通知中にエラーが発生しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}
}
