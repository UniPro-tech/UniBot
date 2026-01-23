package messageComponent

import (
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

// MessageComponentHandler は MessageComponent のハンドラー関数の型
type MessageComponentHandler func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate)

// Handlers はプレフィックスごとのハンドラーを登録するマップ
var Handlers = map[string]MessageComponentHandler{}

// RegisterHandler は指定したプレフィックスに対するハンドラーを登録します
func RegisterHandler(prefix string, handler MessageComponentHandler) {
	Handlers[prefix] = handler
}
