package messageComponent

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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

func init() {
	RegisterHandler("tts_set_voice_select", HandleTTSSetVoice)
	RegisterHandler("tts_set_voice_page", HandleTTSSetVoicePage)
}

// HandleTTSSetVoice は話者選択のセレクトメニューを処理します
func HandleTTSSetVoice(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	values := i.MessageComponentData().Values
	if len(values) == 0 {
		return
	}

	speakerID := values[0]
	memberID := i.Member.User.ID

	memberRepo := repository.NewMemberRepository(ctx.DB)
	if err := memberRepo.Create(memberID); err != nil {
		log.Println("Error creating member:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "メンバー情報の作成に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil {
		log.Println("Error fetching TTS personal setting:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "TTS個人設定の取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if setting == nil {
		defaultSetting := repository.DefaultTTSPersonalSetting
		setting = &defaultSetting
		setting.MemberID = memberID
		setting.SpeakerID = speakerID
		if err := repo.Create(setting); err != nil {
			log.Println("Error creating TTS personal setting:", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "TTS個人設定の作成に失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	} else {
		setting.SpeakerID = speakerID
		if err := repo.Update(setting); err != nil {
			log.Println("Error updating TTS personal setting:", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "TTS個人設定の更新に失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	}

	label := resolveSpeakerLabel(ctx, speakerID)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "話者設定を更新しました",
					Description: "選択した話者: " + label,
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}

// HandleTTSSetVoicePage は話者選択のページ送りを処理します
func HandleTTSSetVoicePage(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID
	parts := strings.Split(customID, ":")
	if len(parts) != 2 {
		return
	}

	pageIndex, err := strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	speakers, err := fetchSpeakers(ctx)
	if err != nil {
		log.Println("Failed to fetch speakers:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "話者情報の取得に失敗しました。",
						Color:       ctx.Config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	pages := buildSpeakerPages(speakers, speakerPageSize)
	if len(pages) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "話者情報が取得できませんでした。",
						Color:       ctx.Config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	currentSpeakerID := getCurrentSpeakerID(ctx, i.Member.User.ID)
	content, components := buildVoiceMessage(pageIndex, pages, currentSpeakerID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content:    content,
			Components: components,
		},
	})
}

func resolveSpeakerLabel(ctx *internal.BotContext, speakerID string) string {
	speakers, err := fetchSpeakers(ctx)
	if err != nil {
		return "ID: " + speakerID
	}

	for _, speaker := range speakers {
		for _, style := range speaker.Styles {
			if fmt.Sprintf("%d", style.ID) == speakerID {
				return fmt.Sprintf("%s / %s", speaker.Name, style.Name)
			}
		}
	}

	return "ID: " + speakerID
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
