package ttsutil

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unibot/internal"
	"unibot/internal/api/voicevox"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

const (
	SpeakerPageSize         = 20
	VoiceSelectCustomID     = "tts_set_voice_select"
	VoicePageCustomIDPrefix = "tts_set_voice_page"
	speakerSelectMax        = 25
	speakerCacheTTL         = 5 * time.Minute
)

type SpeakerPage struct {
	Options []discordgo.SelectMenuOption
}

type speakerCache struct {
	mu       sync.RWMutex
	speakers []voicevox.Speaker
	expires  time.Time
}

var cachedSpeakers speakerCache

func FetchSpeakers(ctx *internal.BotContext) ([]voicevox.Speaker, error) {
	cachedSpeakers.mu.RLock()
	if time.Now().Before(cachedSpeakers.expires) && len(cachedSpeakers.speakers) > 0 {
		speakers := cachedSpeakers.speakers
		cachedSpeakers.mu.RUnlock()
		return speakers, nil
	}
	cachedSpeakers.mu.RUnlock()

	cachedSpeakers.mu.Lock()
	if time.Now().Before(cachedSpeakers.expires) && len(cachedSpeakers.speakers) > 0 {
		speakers := cachedSpeakers.speakers
		cachedSpeakers.mu.Unlock()
		return speakers, nil
	}
	cachedSpeakers.mu.Unlock()

	baseCtx := context.Background()
	if ctx != nil {
		if parentCtx, ok := any(ctx).(context.Context); ok && parentCtx != nil {
			baseCtx = parentCtx
		}
	}

	requestCtx, cancel := context.WithTimeout(baseCtx, 30*time.Second)
	defer cancel()

	speakers, err := ctx.VoiceVox.GetSpeakers(requestCtx)
	if err != nil {
		return nil, err
	}

	cachedSpeakers.mu.Lock()
	cachedSpeakers.speakers = speakers
	cachedSpeakers.expires = time.Now().Add(speakerCacheTTL)
	cachedSpeakers.mu.Unlock()

	return speakers, nil
}

func BuildSpeakerPages(speakers []voicevox.Speaker, perPage int) []SpeakerPage {
	if perPage <= 0 {
		perPage = SpeakerPageSize
	}
	if perPage > speakerSelectMax {
		perPage = speakerSelectMax
	}

	pages := make([]SpeakerPage, 0)
	current := SpeakerPage{Options: make([]discordgo.SelectMenuOption, 0, perPage)}
	flush := func() {
		if len(current.Options) > 0 {
			pages = append(pages, current)
			current = SpeakerPage{Options: make([]discordgo.SelectMenuOption, 0, perPage)}
		}
	}

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

		if len(speakerOptions) == 0 {
			continue
		}

		if len(speakerOptions) > speakerSelectMax {
			flush()
			for start := 0; start < len(speakerOptions); start += speakerSelectMax {
				end := start + speakerSelectMax
				if end > len(speakerOptions) {
					end = len(speakerOptions)
				}
				pages = append(pages, SpeakerPage{Options: speakerOptions[start:end]})
			}
			continue
		}

		if len(speakerOptions) > perPage {
			flush()
			pages = append(pages, SpeakerPage{Options: speakerOptions})
			continue
		}

		if len(current.Options)+len(speakerOptions) > perPage {
			flush()
		}
		current.Options = append(current.Options, speakerOptions...)
	}

	flush()
	return pages
}

func BuildVoiceMessage(pageIndex int, pages []SpeakerPage, currentSpeakerID string) (string, []discordgo.MessageComponent) {
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
					CustomID:    VoiceSelectCustomID,
					Placeholder: "話者を選択してください",
					Options:     pages[pageIndex].Options,
				},
			},
		},
	}

	if maxPage > 1 {
		prevID := fmt.Sprintf("%s:%d", VoicePageCustomIDPrefix, pageIndex-1)
		nextID := fmt.Sprintf("%s:%d", VoicePageCustomIDPrefix, pageIndex+1)
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

func GetCurrentSpeakerID(ctx *internal.BotContext, memberID string) string {
	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil || setting == nil {
		return repository.DefaultTTSPersonalSetting.SpeakerID
	}
	return setting.SpeakerID
}

func ResolveSpeakerLabel(ctx *internal.BotContext, speakerID string) string {
	speakers, err := FetchSpeakers(ctx)
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

func IsSpeakerIDValid(ctx *internal.BotContext, speakerID string) bool {
	speakers, err := FetchSpeakers(ctx)
	if err != nil {
		return false
	}

	for _, speaker := range speakers {
		for _, style := range speaker.Styles {
			if fmt.Sprintf("%d", style.ID) == speakerID {
				return true
			}
		}
	}

	return false
}

func GetInteractionUser(i *discordgo.InteractionCreate) (string, string, string) {
	if i.Member != nil && i.Member.User != nil {
		return i.Member.User.ID, i.Member.DisplayName(), i.Member.AvatarURL("")
	}
	if i.User != nil {
		return i.User.ID, i.User.Username, i.User.AvatarURL("")
	}
	return "", "", ""
}
