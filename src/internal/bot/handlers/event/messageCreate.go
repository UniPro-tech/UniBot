package event_handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/model"
	"unibot/internal/query"
	"unibot/internal/util"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/gorm"
)

// 正規表現パターン
var (
	codeBlockRegex      = regexp.MustCompile("(?s)```(\\w*)\\n.*?```")
	inlineCodeRegex     = regexp.MustCompile("`[^`]*`")
	channelMentionRegex = regexp.MustCompile(`<#(\d+)>`)
	userMentionRegex    = regexp.MustCompile(`<@!?(\d+)>`)
	roleMentionRegex    = regexp.MustCompile(`<@&(\d+)>`)
	customEmojiRegex    = regexp.MustCompile(`<a?:[^:]+:\d+>`) // <:name:id> or <a:name:id>
	unicodeEmojiRegex   = regexp.MustCompile(`[\p{So}\p{Sk}]`) // Unicode絵文字
	urlRegex            = regexp.MustCompile(`https?://[^\s]+`)
	spoilerRegex        = regexp.MustCompile(`\|\|.*?\|\|`)
)

type ExtentionConstant struct {
	Extention []string
	Yomi      string
}

type AttachementTypeList struct {
	ExtentionData       ExtentionConstant
	NumberOfAttachement int
}

// 拡張子一覧
var (
	imageExtensions = ExtentionConstant{
		Extention: []string{".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp"},
		Yomi:      "画像",
	}
	videoExtensions = ExtentionConstant{
		Extention: []string{".mp4", ".mov", ".avi", ".mkv", ".webm"},
		Yomi:      "動画",
	}
	audioExtensions = ExtentionConstant{
		Extention: []string{".mp3", ".wav", ".ogg", ".flac", ".aac"},
		Yomi:      "音声",
	}
	documentExtensions = ExtentionConstant{
		Extention: []string{".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"},
		Yomi:      "文書",
	}
	archiveExtensions = ExtentionConstant{
		Extention: []string{".zip", ".rar", ".7z", ".tar", ".gz"},
		Yomi:      "アーカイブ",
	}
	textExtensions = ExtentionConstant{
		Extention: []string{".txt"},
		Yomi:      "テキスト",
	}
	markdownExtensions = ExtentionConstant{
		Extention: []string{".md", ".markdown"},
		Yomi:      "マークダウン",
	}
	csvExtensions = ExtentionConstant{
		Extention: []string{".csv"},
		Yomi:      "CSV",
	}
	executableExtensions = ExtentionConstant{
		Extention: []string{".exe", ".msi", ".bat", ".sh", ".bin"},
		Yomi:      "実行可能ファイル",
	}
)

var attachmentCategories = []ExtentionConstant{
	imageExtensions,
	videoExtensions,
	audioExtensions,
	documentExtensions,
	archiveExtensions,
	textExtensions,
	markdownExtensions,
	csvExtensions,
	executableExtensions,
}

func MessageCreate(ctx context.Context, bctx *internal.BotContext, e *events.MessageCreate) {
	// Ignore bot itself
	if e.Message.Author.ID == e.Client().ID() {
		return
	}

	// Ignore DM
	if e.GuildID == nil {
		return
	}
	guildId := *e.GuildID

	// ----- TTS -----
	ttsConnectionData, err := query.TtsConnection.Where(query.TtsConnection.GuildID.Eq(int64(guildId))).First()
	if err != nil {
		// TTS 未設定のギルドでは毎メッセージこの分岐に入るため、
		// レコード無しはエラーとして扱わない。
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			slog.WarnContext(ctx, "failed to load tts connection",
				slog.String("guild_id", guildId.String()), slog.Any("err", err))
		}
		return
	}

	if e.Message.Flags&discord.MessageFlagSuppressNotifications != 0 {
		return
	}

	if ttsConnectionData != nil {
		userID := e.Message.Author.ID

		if e.Message.Author.Bot {
			return
		}

		inVC := false
		for vs := range e.Client().Caches.VoiceStates(guildId) {
			if vs.UserID == e.Client().ID() {
				inVC = true
				break
			}
		}

		if inVC {
			var botChannelID *snowflake.ID
			for vs := range e.Client().Caches.VoiceStates(guildId) {
				if vs.UserID == e.Client().ID() {
					botChannelID = vs.ChannelID
					break
				}
			}

			// botChannelID は VoiceState 次第で nil になりうる。
			if int64(e.ChannelID) != ttsConnectionData.ChannelID &&
				(botChannelID == nil || e.ChannelID != *botChannelID) {
				return
			}
		}

		if e.Message.Content == "s" || e.Message.Content == "skip" {
			player := voice.GetManager().GetPlayer(int64(guildId))
			if player != nil {
				player.SkipCurrent()
			}
			return
		}

		memberPreference, err := query.TtsMemberPreference.Where(query.TtsMemberPreference.GuildID.Eq(int64(guildId)), query.TtsMemberPreference.UserID.Eq(int64(userID))).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			slog.WarnContext(ctx, "failed to load tts personal setting", slog.Any("err", err))
			return
		}
		userPreference, err := query.TtsUserPreference.Where(query.TtsUserPreference.UserID.Eq(int64(userID))).First()
		if err != nil && err != gorm.ErrRecordNotFound {
			slog.WarnContext(ctx, "failed to load tts personal setting", slog.Any("err", err))
			return
		}

		// memberPreference を優先し、未設定なら userPreference を使用
		preference := memberPreference

		if preference == nil {
			preference = &model.TtsMemberPreference{
				// 必要なフィールドを userPreference から設定
			}
		}
		if userPreference != nil {
			if preference.SpeakerID == nil {
				preference.SpeakerID = &userPreference.SpeakerID
			}
			if preference.Speed == nil {
				preference.Speed = &userPreference.Speed
			}
		}

		content := SanitizeMessageContent(e.Client(), e.GuildID, e.Message.Content)
		// 辞書を適用
		content = util.ApplyDictionary(ctx, bctx.DB, int64(*e.GuildID), content)

		// 切り詰め
		content = TruncateForTTS(content, 250)

		// 添付ファイル一覧を取得
		attachmentCounts := map[string]*AttachementTypeList{}

		for _, attachment := range e.Message.Attachments {
			attachmentType := DetectAttachmentType(attachment.Filename)

			if data, exists := attachmentCounts[attachmentType.Yomi]; exists {
				data.NumberOfAttachement++
			} else {
				attachmentCounts[attachmentType.Yomi] = &AttachementTypeList{
					ExtentionData:       attachmentType,
					NumberOfAttachement: 1,
				}
			}
		}

		// 添付ファイルの説明を生成
		if len(attachmentCounts) > 0 {
			var attachmentDescriptions []string
			for _, data := range attachmentCounts {
				desc := fmt.Sprintf("%sが%dつ", data.ExtentionData.Yomi, data.NumberOfAttachement)
				attachmentDescriptions = append(attachmentDescriptions, desc)
			}
			content += "、" + strings.Join(attachmentDescriptions, "、") + "添付されています。"
		}

		vcConn := e.Client().VoiceManager.GetConn(guildId)

		vp := voice.GetManager().GetOrCreatePlayer(
			int64(guildId),
			ttsConnectionData.ChannelID,
			vcConn,
			bctx,
		)

		vp.EnqueueText(voice.QueueItem{
			// 読み上げは非同期に処理されるため、cancel を切って trace のみ引き継ぐ。
			Ctx:                context.WithoutCancel(ctx),
			Text:               content,
			ResolvedPreference: *preference,
		})
	}
}

// メッセージ内容をサニタイズする関数
func SanitizeMessageContent(client *bot.Client, guildID *snowflake.ID, content string) string {
	// コードブロック置換
	content = codeBlockRegex.ReplaceAllStringFunc(content, func(block string) string {
		matches := codeBlockRegex.FindStringSubmatch(block)
		lang := ""
		if len(matches) > 1 {
			lang = matches[1]
		}
		if lang != "" {
			return "、(" + lang + "のコードブロック省略)、"
		}
		return "、(コードブロック省略)、"
	})

	// インラインコード置換
	content = inlineCodeRegex.ReplaceAllString(content, "、(インラインコード省略)、")

	// チャンネルメンション置換
	content = channelMentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := channelMentionRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		channelID := matches[1]
		channel, ok := client.Caches.Channel(snowflake.MustParse(channelID))
		if !ok {
			return match
		}
		return "#" + channel.Name()
	})

	// ユーザーメンション置換
	content = userMentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := userMentionRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		userIDStr := matches[1]
		uID, err := snowflake.Parse(userIDStr)
		if err != nil {
			return match
		}

		// 1. まずキャッシュを探す
		if member, ok := client.Caches.Member(*guildID, uID); ok {
			return "@" + member.EffectiveName()
		}

		// 2. キャッシュにない場合、REST APIで取得すると重いので
		// 一旦そのままにするか、Userキャッシュだけでも探す
		if user, ok := client.Caches.Member(*guildID, uID); ok {
			return "@" + user.EffectiveName()
		}

		// 3. どうしても名前が取れない場合は「不明なユーザー」等にする
		return "@不明なユーザー"
	})

	// ロールメンション置換
	content = roleMentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := roleMentionRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		roleID := matches[1]
		role, ok := client.Caches.Role(*guildID, snowflake.MustParse(roleID))
		if ok {
			return "@" + role.Name
		}
		return match
	})

	// カスタム絵文字置換
	content = customEmojiRegex.ReplaceAllString(content, "、(絵文字)、")

	// Unicode絵文字置換
	content = unicodeEmojiRegex.ReplaceAllString(content, "、(絵文字)、")

	// URL置換
	content = urlRegex.ReplaceAllString(content, "、(リンク省略)、")

	// スポイラー置換
	content = spoilerRegex.ReplaceAllString(content, "、(スポイラー)、")

	return content
}

// TTS用にメッセージを切り詰める関数
func TruncateForTTS(content string, maxLen int) string {
	if len([]rune(content)) <= maxLen {
		return content
	}

	runes := []rune(content)
	cut := maxLen

	// 「、」または「。」で区切れる位置を探す
	for i := maxLen - 1; i >= 0; i-- {
		if runes[i] == '、' || runes[i] == '。' {
			cut = i + 1 // ここで切る
			break
		}
	}

	return string(runes[:cut]) + " 、以下省略"
}

func DetectAttachmentType(filename string) ExtentionConstant {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, category := range attachmentCategories {
		if slices.Contains(category.Extention, ext) {
			return category
		}
	}

	return ExtentionConstant{
		Yomi: "その他",
	}
}
