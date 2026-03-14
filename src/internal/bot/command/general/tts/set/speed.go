package set

import (
	"fmt"
	"log"
	"unibot/internal"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func LoadSpeedCommandContext() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "speed",
		Description: "TTSの再生速度を設定します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "value",
				Description: "再生速度（50-200、100=通常速度）",
				Required:    true,
				MinValue:    floatPtr(50),
				MaxValue:    200,
			},
		},
	}
}

func Speed(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	if len(options) == 0 {
		safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "コマンド引数の解析に失敗しました。", ctx.Config.Colors.Error))
		return
	}

	subCommandGroup := options[0]
	if len(subCommandGroup.Options) == 0 {
		safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "コマンド引数の解析に失敗しました。", ctx.Config.Colors.Error))
		return
	}

	subCommand := subCommandGroup.Options[0]
	if len(subCommand.Options) == 0 {
		safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "コマンド引数の解析に失敗しました。", ctx.Config.Colors.Error))
		return
	}

	speed := subCommand.Options[0].IntValue()

	HandleSpeedCommand(s, i, ctx, speed)
}

func HandleSpeedCommand(s *discordgo.Session, i *discordgo.InteractionCreate, ctx *internal.BotContext, speed int64) {
	memberID := i.Member.User.ID
	memberRepo := repository.NewMemberRepository(ctx.DB)
	if err := memberRepo.Create(memberID); err != nil {
		log.Println("Error creating member:", err)
		safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "メンバー情報の作成に失敗しました。", ctx.Config.Colors.Error))
		return
	}

	repo := repository.NewTTSPersonalSettingRepository(ctx.DB)
	setting, err := repo.GetByMember(memberID)
	if err != nil {
		log.Println("Error fetching TTS personal setting:", err)
		safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "TTS個人設定の取得に失敗しました。", ctx.Config.Colors.Error))
		return
	}
	if setting == nil {
		defaultSetting := repository.DefaultTTSPersonalSetting
		setting = &defaultSetting
		setting.MemberID = memberID
		setting.SpeakerSpeed = speed
		err = repo.Create(setting)
		if err != nil {
			log.Println("Error creating TTS personal setting:", err)
			safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "TTS個人設定の作成に失敗しました。", ctx.Config.Colors.Error))
			return
		}
	} else {
		setting.SpeakerSpeed = speed
		err = repo.Update(setting)
		if err != nil {
			log.Println("Error updating TTS personal setting:", err)
			safeEditSpeedResponse(s, i, buildSpeedEmbed("エラー", "TTS個人設定の更新に失敗しました。", ctx.Config.Colors.Error))
			return
		}
	}
	safeEditSpeedResponse(s, i, buildSpeedEmbed("TTS再生速度設定", "TTSの再生速度を設定しました: "+formatSpeed(speed), ctx.Config.Colors.Success))
}

func buildSpeedEmbed(title, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}

func safeEditSpeedResponse(s *discordgo.Session, i *discordgo.InteractionCreate, embed *discordgo.MessageEmbed) {
	_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
		Flags:  discordgo.MessageFlagsEphemeral,
	})
	if err != nil {
		log.Println("Failed to edit deferred interaction (speed):", err)
	}
}

// formatSpeed はSpeedScale値(100 = 1.0倍)を読みやすい形式に変換する
func formatSpeed(speed int64) string {
	return fmt.Sprintf("%.2f倍速", float64(speed)/100.0)
}

// floatPtr はfloat64のポインタを返すヘルパー関数
func floatPtr(f float64) *float64 {
	return &f
}
