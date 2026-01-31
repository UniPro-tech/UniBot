package schedule

import (
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadRemoveCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "remove",
		Description: "予約投稿を削除します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "jobid",
				Description: "削除する予約投稿のジョブID",
				Required:    true,
			},
		},
	}
}

func Remove(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	options := i.ApplicationCommandData().Options[0].Options

	var jobID string
	for _, opt := range options {
		if opt.Name == "jobid" {
			jobID = opt.StringValue()
		}
	}

	if jobID == "" {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "削除するジョブIDを指定してください。",
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

	repo := repository.NewScheduleSettingRepository(ctx.DB)
	setting, err := repo.GetByID(jobID)
	if err != nil || setting == nil {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "指定されたジョブIDが見つかりません。",
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

	if setting.GuildID != i.GuildID {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "このジョブを削除する権限がありません。",
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

	err = repo.DeleteByID(jobID)
	if err != nil {
		_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "スケジュールの削除に失敗しました。",
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

	successEmbed := &discordgo.MessageEmbed{
		Title:       "予約投稿を削除しました",
		Description: "指定された予約投稿を削除しました。",
		Color:       config.Colors.Success,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	_ = RespondEdit(s, i, &discordgo.InteractionResponseData{
		Embeds: []*discordgo.MessageEmbed{successEmbed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
}
