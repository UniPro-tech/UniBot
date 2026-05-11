package event_handlers

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"
	"unibot/internal/util"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
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

func MessageCreate(ctx *internal.BotContext, e *events.MessageCreate) {
	// Ignore bot itself
	if e.Message.Author.ID == e.Client().ID() {
		return
	}

	// Ignore DM
	if e.GuildID == nil {
		return
	}

	// ----- TTS -----

	repo := repository.NewTTSConnectionRepository(ctx.DB)

	ttsConnectionData, err := repo.GetByGuildID(e.GuildID.String())
	if err != nil {
		log.Println(err)
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
		for vs := range e.Client().Caches.VoiceStates(*e.GuildID) {
			if vs.UserID == e.Client().ID() {
				inVC = true
				break
			}
		}

		if inVC {
			var botChannelID *snowflake.ID
			for vs := range e.Client().Caches.VoiceStates(*e.GuildID) {
				if vs.UserID == e.Client().ID() {
					botChannelID = vs.ChannelID
					break
				}
			}

			if e.ChannelID.String() != ttsConnectionData.ChannelID &&
				e.ChannelID != *botChannelID {
				return
			}
		}

		if e.Message.Content == "s" || e.Message.Content == "skip" {
			player := voice.GetManager().Get(e.GuildID.String())
			if player != nil {
				player.SkipCurrent()
			}
			return
		}

		personalSetting, err := repository.NewTTSPersonalSettingRepository(ctx.DB).GetByMember(userID.String())
		if err != nil {
			log.Println(err)
			return
		}
		if personalSetting == nil {
			personalSetting = &repository.DefaultTTSPersonalSetting
		}
		content := SanitizeMessageContent(e.Client(), e.GuildID, e.Message.Content)

		// 辞書を適用
		content = util.ApplyDictionary(ctx.DB, e.GuildID.String(), content)

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

		vcConn := e.Client().VoiceManager.GetConn(*e.GuildID)

		vp := voice.GetManager().GetOrCreate(
			e.GuildID.String(),
			ttsConnectionData.ChannelID,
			vcConn,
			ctx,
		)

		vp.EnqueueText(voice.QueueItem{
			Text:    content,
			Setting: *personalSetting,
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
