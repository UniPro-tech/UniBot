package handler

import (
	"log"
	"regexp"
	"strings"
	"unibot/internal"
	"unibot/internal/bot/voice"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(ctx *internal.BotContext) func(s *discordgo.Session, r *discordgo.MessageCreate) {
	return func(s *discordgo.Session, r *discordgo.MessageCreate) {

		// Ignore bot itself
		if r.Author.ID == s.State.User.ID {
			return
		}

		// Ignore DM
		if r.GuildID == "" {
			return
		}

		// ----- TTS -----

		repo := repository.NewTTSConnectionRepository(ctx.DB)

		ttsConnectionData, err := repo.GetByGuildID(r.GuildID)
		if err != nil {
			log.Println(err)
			return
		}

		if r.Flags&discordgo.MessageFlagsSuppressNotifications != 0 {
			return
		}

		if ttsConnectionData != nil {
			userID := r.Author.ID

			if r.Author.Bot {
				return
			}

			if r.ChannelID != ttsConnectionData.ChannelID && r.ChannelID != s.VoiceConnections[r.GuildID].ChannelID {
				return
			}

			if r.Content == "s" || r.Content == "skip" {
				player := voice.GetManager().Get(r.GuildID)
				if player != nil {
					player.SkipCurrent()
				}
				return
			}

			personalSetting, err := repository.NewTTSPersonalSettingRepository(ctx.DB).GetByMember(userID)
			if err != nil {
				log.Println(err)
				return
			}
			if personalSetting == nil {
				personalSetting = &repository.DefaultTTSPersonalSetting
			}
			content := SanitizeMessageContent(s, r.GuildID, r.Content)

			// 辞書を適用
			content = ApplyDictionary(ctx, r.GuildID, content)

			content = TruncateForTTS(content, 250)

			vp := voice.GetManager().GetOrCreate(r.GuildID, s.VoiceConnections[r.GuildID], ctx)

			vp.EnqueueText(voice.QueueItem{
				Text:    content,
				Setting: *personalSetting,
			})
		}
	}
}

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
)

// メッセージ内容をサニタイズする関数
func SanitizeMessageContent(s *discordgo.Session, guildID, content string) string {
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
		channel, err := s.State.Channel(channelID)
		if err != nil {
			channel, err = s.Channel(channelID)
			if err != nil {
				return match
			}
		}
		return "#" + channel.Name
	})

	// ユーザーメンション置換
	content = userMentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := userMentionRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		userID := matches[1]
		user, err := s.User(userID)
		if err != nil {
			user, err = s.User(userID)
			if err != nil {
				return match
			}
		}
		return "@" + user.Username
	})

	// ロールメンション置換
	content = roleMentionRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := roleMentionRegex.FindStringSubmatch(match)
		if len(matches) < 2 {
			return match
		}
		roleID := matches[1]
		guild, err := s.State.Guild(guildID)
		if err != nil {
			guild, err = s.Guild(guildID)
			if err != nil {
				return match
			}
		}
		for _, role := range guild.Roles {
			if role.ID == roleID {
				return "@" + role.Name
			}
		}
		return match
	})

	// カスタム絵文字置換
	content = customEmojiRegex.ReplaceAllString(content, "、(絵文字)、")

	// Unicode絵文字置換
	content = unicodeEmojiRegex.ReplaceAllString(content, "、(絵文字)、")

	// URL置換
	content = urlRegex.ReplaceAllString(content, "、(リンク省略)、")

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

// 辞書を適用する関数
func ApplyDictionary(ctx *internal.BotContext, guildID, content string) string {
	repo := repository.NewTTSDictionaryRepository(ctx.DB)

	entries, err := repo.ListByGuild(guildID)
	if err != nil {
		log.Println("辞書の取得に失敗しました:", err)
		return content
	}

	for _, entry := range entries {
		if entry.CaseSensitive {
			// 大文字小文字を区別して置換
			content = strings.ReplaceAll(content, entry.Word, entry.Definition)
		} else {
			// 大文字小文字を区別せずに置換
			content = replaceIgnoreCase(content, entry.Word, entry.Definition)
		}
	}

	return content
}

// 大文字小文字を無視して置換する関数
func replaceIgnoreCase(input, old, new string) string {
	lowerInput := strings.ToLower(input)
	lowerOld := strings.ToLower(old)

	var result strings.Builder
	i := 0

	for i < len(input) {
		idx := strings.Index(lowerInput[i:], lowerOld)
		if idx == -1 {
			result.WriteString(input[i:])
			break
		}

		result.WriteString(input[i : i+idx])
		result.WriteString(new)
		i += idx + len(old)
	}

	return result.String()
}
