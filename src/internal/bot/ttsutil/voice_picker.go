package ttsutil

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unibot/internal"
	"unibot/internal/api/voicevox"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
)

const (
	SpeakerPageSize         = 20
	VoiceSelectCustomID     = "tts_set_voice_select"
	VoicePageCustomIDPrefix = "tts_set_voice_page"
	speakerSelectMax        = 25
	speakerCacheTTL         = 5 * time.Minute
)

type SpeakerPage struct {
	Options []discord.StringSelectMenuOption
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

	requestCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	speakers, err := ctx.VoiceVox.GetSpeakers(requestCtx)
	if err != nil {
		cachedSpeakers.mu.Unlock()
		return nil, err
	}

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
	current := SpeakerPage{Options: make([]discord.StringSelectMenuOption, 0, perPage)}
	flush := func() {
		if len(current.Options) > 0 {
			pages = append(pages, current)
			current = SpeakerPage{Options: make([]discord.StringSelectMenuOption, 0, perPage)}
		}
	}

	for _, speaker := range speakers {
		speakerOptions := make([]discord.StringSelectMenuOption, 0, len(speaker.Styles))
		for _, style := range speaker.Styles {
			label := fmt.Sprintf("%s / %s", speaker.Name, style.Name)
			speakerOptions = append(speakerOptions, discord.StringSelectMenuOption{
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

func BuildVoiceMessage(pageIndex int, pages []SpeakerPage, currentSpeakerID string) (string, []discord.LayoutComponent) {
	maxPage := len(pages)
	if maxPage == 0 {
		return "話者情報が取得できませんでした。", []discord.LayoutComponent{}
	}
	if pageIndex < 0 {
		pageIndex = 0
	}
	if pageIndex >= maxPage {
		pageIndex = maxPage - 1
	}

	content := fmt.Sprintf("話者を選択してください。\n現在の話者ID: %s\nページ %d/%d", currentSpeakerID, pageIndex+1, maxPage)

	components := []discord.LayoutComponent{
		discord.ActionRowComponent{
			Components: []discord.InteractiveComponent{
				discord.NewStringSelectMenu(
					VoiceSelectCustomID,
					"話者を選択してください",
				).SetOptions(pages[pageIndex].Options...),
			},
		},
	}

	if maxPage > 1 {
		prevID := fmt.Sprintf("%s:%d", VoicePageCustomIDPrefix, pageIndex-1)
		nextID := fmt.Sprintf("%s:%d", VoicePageCustomIDPrefix, pageIndex+1)
		components = append(components, discord.ActionRowComponent{
			Components: []discord.InteractiveComponent{
				discord.ButtonComponent{
					CustomID: prevID,
					Label:    "前へ",
					Style:    discord.ButtonStyleSecondary,
					Disabled: pageIndex == 0,
				},
				discord.ButtonComponent{
					CustomID: nextID,
					Label:    "次へ",
					Style:    discord.ButtonStyleSecondary,
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
