package ttsSet

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
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
		},
	}
}

func Speed(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		speedValue := data.Options["speed"].Value

		var speed int
		err := json.Unmarshal(speedValue, &speed)
		if err != nil {
			log.Print(err)
		}

		return handleSpeedCommand(e, ctx, speed)
	}
}

func handleSpeedCommand(e *handler.CommandEvent, ctx *internal.BotContext, speed int) error {
	requester := e.User()

	memberID := requester.ID.String()
	if speed < MinSpeakerSpeed || speed > MaxSpeakerSpeed {
		responseEmbed := buildSpeedEmbed("エラー", fmt.Sprintf("再生速度は%d〜%dの範囲で指定してください。", MinSpeakerSpeed, MaxSpeakerSpeed), ctx.Config.Colors.Error, &requester)
		_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
		return err
	}

	memberRepo := repository.NewMemberRepository(ctx.DB)
	if err := memberRepo.Create(memberID); err != nil {
		log.Println("Error creating member:", err)
		responseEmbed := buildSpeedEmbed("エラー", "メンバー情報の作成に失敗しました。", ctx.Config.Colors.Error, &requester)
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
		return err
	}

	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil {
		log.Println("Error fetching TTS personal setting:", err)
		responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の取得に失敗しました。", ctx.Config.Colors.Error, &requester)
		_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
		return err
	}
	if setting == nil {
		defaultSetting := repository.DefaultTTSPersonalSetting
		setting = &defaultSetting
		setting.MemberID = memberID
		setting.SpeakerSpeed = speed
		err = repo.Create(setting)
		if err != nil {
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の作成に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
	} else {
		setting.SpeakerSpeed = speed
		err = repo.Update(setting)
		if err != nil {
			responseEmbed := buildSpeedEmbed("エラー", "TTS個人設定の更新に失敗しました。", ctx.Config.Colors.Error, &requester)
			_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
			return err
		}
	}
	responseEmbed := buildSpeedEmbed("TTS再生速度設定", "TTSの再生速度を設定しました: "+formatSpeed(speed), ctx.Config.Colors.Success, &requester)
	_, err = e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(*responseEmbed))
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
