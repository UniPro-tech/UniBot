package handler

import (
	"log"
	"strings"
	"unibot/internal"
	"unibot/internal/bot/command"
	"unibot/internal/bot/messageComponent"

	"github.com/disgoorg/disgo/events"
)

// IdisgoのApplicationCommandInteractionイベントを処理する
func InteractionCreate(ctx *internal.BotContext) func(e *events.ApplicationCommandInteractionCreate) {
	return func(e *events.ApplicationCommandInteractionCreate) {
		handleApplicationCommand(ctx, e)
	}
}

// ボタンやセレクトメニューのイベントを処理
func ComponentInteractionCreate(ctx *internal.BotContext) func(e *events.ComponentInteractionCreate) {
	return func(e *events.ComponentInteractionCreate) {
		handleMessageComponent(ctx, e)
	}
}

func handleApplicationCommand(ctx *internal.BotContext, e *events.ApplicationCommandInteractionCreate) {
	name := e.Data.CommandName()

	// 1. エフェメラル（自分にだけ見える）設定の判定
	isEphemeral := false
	if entry, ok := command.Handlers[name]; (ok && entry.Ephemeral) || isTtsSetCommand(e) {
		isEphemeral = true
	}

	// 2. Defer (考え中... の状態にする)
	if err := e.DeferCreateMessage(isEphemeral); err != nil {
		log.Println("Failed to respond interaction:", err)
	}

	// 3. 実際の処理を実行
	if entry, ok := command.Handlers[name]; ok {
		entry.Handler(ctx, e)
	}
}

func isTtsSetCommand(e *events.ApplicationCommandInteractionCreate) bool {
	data := e.SlashCommandInteractionData()

	if data.CommandName() != "tts" {
		return false
	}

	// .SubCommandGroupName は *string 型のフィールド
	// 1. nil チェック (その階層があるか)
	// 2. デリファレンスして値の比較
	if data.SubCommandGroupName != nil && *data.SubCommandGroupName == "set" {
		return true
	}

	// グループがない直下のサブコマンドの場合
	if data.SubCommandName != nil && *data.SubCommandName == "set" {
		return true
	}

	return false
}

func handleMessageComponent(ctx *internal.BotContext, e *events.ComponentInteractionCreate) {
	customID := e.Data.CustomID()

	for prefix, handler := range messageComponent.Handlers {
		if strings.HasPrefix(customID, prefix) {
			handler(ctx, e)
			return
		}
	}
}
