package messageComponent

import (
	"strings"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

// MessageComponentHandler は MessageComponent のハンドラー関数の型
type MessageComponentHandler func(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate)

// handlers はプレフィックスごとのハンドラーを登録するマップ
var handlers = map[string]MessageComponentHandler{}

// RegisterHandler は指定したプレフィックスに対するハンドラーを登録します
func RegisterHandler(prefix string, handler MessageComponentHandler) {
	handlers[prefix] = handler
}

// Handle はカスタムIDに基づいて適切なハンドラーを呼び出します
func Handle(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	for prefix, handler := range handlers {
		if strings.HasPrefix(customID, prefix) {
			handler(ctx, s, i)
			return
		}
	}
}
