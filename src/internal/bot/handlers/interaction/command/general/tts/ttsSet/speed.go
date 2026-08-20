package ttsSet

import (
	"fmt"
	"log"
	"log/slog"
	"time"
	"unibot/internal"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"gorm.io/gorm"
)

const (
	MinSpeakerSpeed int = 50
	MaxSpeakerSpeed int = 200
)

func LoadSpeedCommandContext() discord.ApplicationCommandOptionSubCommand {
	return discord.ApplicationCommandOptionSubCommand{
		Name:        "speed",
		Description: "TTSの再生速度を設定します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionInt{
				Name:        "speed",
				Description: "再生速度（50-200、100=通常速度）",
				Required:    true,
				MinValue:    intPtr(MinSpeakerSpeed),
				MaxValue:    intPtr(MaxSpeakerSpeed),
			},
			discord.ApplicationCommandOptionBool{
				Name:        "global",
				Description: "グローバルに設定するか",
				Required:    false,
			},
		},
	}
}

func Speed(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		speed := data.Options["speed"].Int()

		global := false
		if value, ok := data.Options["global"]; ok {
			global = value.Bool()
		}

		return handleSpeedCommand(e, ctx, speed, global)
	}
}

func handleSpeedCommand(e *handler.CommandEvent, ctx *internal.BotContext, speed int, isGlobal bool) error {
	requester := e.User()

	memberID := requester.ID
	if speed < MinSpeakerSpeed || speed > MaxSpeakerSpeed {
		responseEmbed := buildSpeedEmbed("エラー", fmt.Sprintf("再生速度は%d〜%dの範囲で指定してください。", MinSpeakerSpeed, MaxSpeakerSpeed), ctx.Config.Colors.Error, &requester)
		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
		return err
	}

	if isGlobal {
		setting, err := query.TtsUserPreference.Where(query.TtsUserPreference.UserID.Eq(int64(memberID))).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			slog.ErrorContext(e.Ctx, "failed to fetch tts personal setting", slog.Any("err", err))
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の取得に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
		setting.Speed = int32(speed)
		err = query.TtsUserPreference.Save(setting)
		if err != nil {
			log.Println("Error fetching TTS personal setting:", err)
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の保存に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
	} else {
		guildID := e.GuildID()
		setting, err := query.TtsMemberPreference.Where(query.TtsMemberPreference.UserID.Eq(int64(memberID)), query.TtsMemberPreference.GuildID.Eq(int64(*guildID))).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Println("Error fetching TTS personal setting:", err)
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の取得に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
		setting.Speed = int32(speed)
		err = query.TtsMemberPreference.Save(setting)
		if err != nil {
			log.Println("Error fetching TTS personal setting:", err)
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の保存に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
	}

	responseEmbed := buildSpeedEmbed("TTS再生速度設定", "TTSの再生速度を設定しました: "+formatSpeed(speed), ctx.Config.Colors.Success, &requester)
	_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
	return err
}

func buildSpeedEmbed(title, description string, color int, requester *discord.User) *discord.Embed {
	embed := &discord.Embed{
		Title:       title,
		Description: description,
		Color:       color,
		Footer: &discord.EmbedFooter{
			Text:    fmt.Sprintf("Requested by %s", requester.Username),
			IconURL: *requester.AvatarURL(),
		},
		Timestamp: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	}

	return embed
}

// formatSpeed はSpeedScale値(100 = 1.0倍)を読みやすい形式に変換する
func formatSpeed(speed int) string {
	return fmt.Sprintf("%.2f倍速", float64(speed)/100.0)
}

func intPtr(f int) *int {
	return &f
}
