package set

import (
	"context"
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/api/voicevox"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

const speakerPageSize = 20

type speakerPage struct {
	Options []discordgo.SelectMenuOption
}

func LoadVoiceCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "voice",
		Description: "読み上げの話者を設定します",
	}
}

func Voice(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
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
						Description: "話者情報の取得に失敗しました。",
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

	speakers, err := fetchSpeakers(ctx)
	if err != nil {
		log.Println("Failed to fetch speakers:", err)
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "話者情報の取得に失敗しました。",
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

	pages := buildSpeakerPages(speakers, speakerPageSize)
	if len(pages) == 0 {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "話者情報が取得できませんでした。",
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

	currentSpeakerID := getCurrentSpeakerID(ctx, i.Member.User.ID)
	content, components := buildVoiceMessage(0, pages, currentSpeakerID)

	s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content:    &content,
		Components: &components,
	})
}

func fetchSpeakers(ctx *internal.BotContext) ([]voicevox.Speaker, error) {
	requestCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return ctx.VoiceVox.GetSpeakers(requestCtx)
}

func buildSpeakerPages(speakers []voicevox.Speaker, perPage int) []speakerPage {
	pages := make([]speakerPage, 0)
	current := speakerPage{Options: []discordgo.SelectMenuOption{}}

	for _, speaker := range speakers {
		speakerOptions := make([]discordgo.SelectMenuOption, 0, len(speaker.Styles))
		for _, style := range speaker.Styles {
			label := fmt.Sprintf("%s / %s", speaker.Name, style.Name)
			speakerOptions = append(speakerOptions, discordgo.SelectMenuOption{
				Label:       label,
				Value:       fmt.Sprintf("%d", style.ID),
				Description: fmt.Sprintf("ID: %d", style.ID),
			})
		}

		if len(current.Options) > 0 && len(current.Options)+len(speakerOptions) > perPage {
			pages = append(pages, current)
			current = speakerPage{Options: []discordgo.SelectMenuOption{}}
		}
		current.Options = append(current.Options, speakerOptions...)
	}

	if len(current.Options) > 0 {
		pages = append(pages, current)
	}

	return pages
}

func buildVoiceMessage(pageIndex int, pages []speakerPage, currentSpeakerID string) (string, []discordgo.MessageComponent) {
	maxPage := len(pages)
	if maxPage == 0 {
		return "話者情報が取得できませんでした。", []discordgo.MessageComponent{}
	}
	if pageIndex < 0 {
		pageIndex = 0
	}
	if pageIndex >= maxPage {
		pageIndex = maxPage - 1
	}

	content := fmt.Sprintf("話者を選択してください。\n現在の話者ID: %s\nページ %d/%d", currentSpeakerID, pageIndex+1, maxPage)

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "tts_set_voice_select",
					Placeholder: "話者を選択してください",
					Options:     pages[pageIndex].Options,
				},
			},
		},
	}

	if maxPage > 1 {
		prevID := fmt.Sprintf("tts_set_voice_page:%d", pageIndex-1)
		nextID := fmt.Sprintf("tts_set_voice_page:%d", pageIndex+1)
		components = append(components, discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					CustomID: prevID,
					Label:    "前へ",
					Style:    discordgo.SecondaryButton,
					Disabled: pageIndex == 0,
				},
				discordgo.Button{
					CustomID: nextID,
					Label:    "次へ",
					Style:    discordgo.SecondaryButton,
					Disabled: pageIndex >= maxPage-1,
				},
			},
		})
	}

	return content, components
}

func getCurrentSpeakerID(ctx *internal.BotContext, memberID string) string {
	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil || setting == nil {
		return repository.DefaultTTSPersonalSetting.SpeakerID
	}
	return setting.SpeakerID
}
