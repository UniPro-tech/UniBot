package interaction_handler

import (
	"log/slog"
	"slices"
	"time"

	"unibot/internal"
	"unibot/internal/bot/handlers/interaction/command/admin/maintenance"
	"unibot/internal/bot/handlers/interaction/command/general"
	"unibot/internal/bot/handlers/interaction/command/general/rss"
	"unibot/internal/bot/handlers/interaction/command/general/tts"
	"unibot/internal/bot/handlers/interaction/command/general/tts/dict"
	"unibot/internal/bot/handlers/interaction/command/general/tts/ttsSet"
	"unibot/internal/bot/handlers/interaction/messageComponent"
	"unibot/internal/query"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

func RegistHandler(r *handler.Mux, ctxData *internal.BotContext) {
	r.Route("/ping", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, false, false))
		r.SlashCommand("/", general.Ping(ctxData))
	})
	r.Route("/about", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, false, false))
		r.SlashCommand("/", general.About(ctxData))
	})
	r.Route("/help", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, false, false))
		r.SlashCommand("/", general.Help(ctxData))
	})
	r.Route("/colorcode", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, false, false))
		r.SlashCommand("/", general.ColorCode(ctxData))
	})
	r.Route("/tts", func(r handler.Router) {
		r.Route("/join", func(r handler.Router) {
			r.Use(DeferReplyMiddleware(ctxData, false, false))
			r.SlashCommand("/", tts.Join(ctxData))
		})
		r.Route("/leave", func(r handler.Router) {
			r.Use(DeferReplyMiddleware(ctxData, false, false))
			r.SlashCommand("/", tts.Leave(ctxData))
		})
		r.Route("/skip", func(r handler.Router) {
			r.Use(DeferReplyMiddleware(ctxData, false, false))
			r.SlashCommand("/", tts.Skip(ctxData))
		})
		r.Route("/set", func(r handler.Router) {
			r.Use(DeferReplyMiddleware(ctxData, true, false))
			r.SlashCommand("/speed", ttsSet.Speed(ctxData))
			r.SlashCommand("/voice", ttsSet.Voice(ctxData))
		})
		r.Route("/dict", func(r handler.Router) {
			r.Use(DeferReplyMiddleware(ctxData, true, false))
			r.SlashCommand("/add", dict.Add(ctxData))
			r.SlashCommand("/list", dict.List(ctxData))
			r.SlashCommand("/remove", dict.Remove(ctxData))
		})
	})
	r.Route("/maintenance", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, true, false))
		r.Use(AdminOnlyMiddleware(ctxData))
		r.SlashCommand("/status/set", maintenance.StatusSetHandler(ctxData))
		r.SlashCommand("/status/reset", maintenance.StatusResetHandler(ctxData))
		r.SlashCommand("/shutdown", maintenance.Shutdown(ctxData))
	})
	r.Route("/rss", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, true, false))
		r.SlashCommand("/subscribe", rss.Subscribe(ctxData))
	})
	// action row
	// select menu
	r.Route("/tts_dict_remove", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, true, true))
		r.SelectMenuComponent("/", messageComponent.HandleTTSDictRemove(ctxData))
	})
	r.Route("/tts_set_voice_select", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, true, true))
		r.SelectMenuComponent("/", messageComponent.HandleTTSSetVoice(ctxData))
	})
	// button
	r.Route("/tts_set_voice_page/{pageIndex}", func(r handler.Router) {
		r.Use(DeferReplyMiddleware(ctxData, true, true))
		r.ButtonComponent("/", messageComponent.HandleTTSSetVoicePage(ctxData))
	})
}

// IsOwner は管理者ロールを持つメンバーかどうかを判定する。
//
// 設定はリクエストごとに読み直さず、呼び出し側から渡してもらう。
func IsOwner(config *internal.Config, member discord.Member) bool {
	adminRoleID, err := snowflake.Parse(config.AdminRoleID)
	if err != nil {
		slog.Error("invalid admin role id, denying admin access",
			slog.String("value", config.AdminRoleID), slog.Any("err", err))
		return false
	}
	return slices.Contains(member.RoleIDs, adminRoleID)
}

func AdminOnlyMiddleware(ctx *internal.BotContext) func(next handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			config := ctx.Config
			member := e.Member()
			if member == nil || !IsOwner(config, member.Member) {
				errorEmbed := discord.Embed{
					Title:       "権限エラー",
					Description: "権限がありません。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						// Nick は未設定のことがあるため EffectiveName を使う。
						Text:    "Requested by " + e.User().EffectiveName(),
						IconURL: e.User().EffectiveAvatarURL(),
					},
					Timestamp: func() *time.Time {
						t := time.Now()
						return &t
					}(),
				}
				_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(errorEmbed).WithEphemeral(true))
				return err
			}

			return next(e)
		}
	}
}

func DeferReplyMiddleware(ctx *internal.BotContext, ephemeral bool, update bool) func(next handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		if !update {
			return func(e *handler.InteractionEvent) error {
				if err := deferInteraction(e, func() error { return e.DeferCreateMessage(ephemeral) }); err != nil {
					return err
				}
				return next(e)
			}
		} else {
			return func(e *handler.InteractionEvent) error {
				if err := deferInteraction(e, func() error { return e.DeferUpdateMessage() }); err != nil {
					return err
				}
				return next(e)
			}
		}
	}
}

// deferInteraction は応答を保留し、その結果を記録する。
//
// 保留に失敗すると以降の followup もすべて失敗するため、エラーを返して
// ハンドラ本体を実行させない。呼び出し元が返した error は
// ObservabilityMiddleware が受け取り、未応答の状態からユーザーへ通知する。
// 成功した場合はエラー応答の送り方を切り替えられるよう context に印を付ける。
func deferInteraction(e *handler.InteractionEvent, deferFn func() error) error {
	if err := deferFn(); err != nil {
		slog.WarnContext(e.Ctx, "failed to defer interaction response", slog.Any("err", err))
		return err
	}
	e.Ctx = markDeferred(e.Ctx)
	return nil
}

func CreateMasterRecordMiddleware(next handler.Handler) handler.Handler {
	return func(e *handler.InteractionEvent) error {
		guild, ok := e.Guild()
		channel := e.Channel()
		if ok {
			if _, err := query.Guild.Where(query.Guild.ID.Eq(int64(guild.ID))).FirstOrCreate(); err != nil {
				return err
			}
			if _, err := query.Channel.Where(query.Channel.ID.Eq(int64(channel.ID())), query.Channel.GuildID.Eq(int64(guild.ID))).FirstOrCreate(); err != nil {
				return err
			}
		} else {
			if _, err := query.Channel.Where(query.Channel.ID.Eq(int64(channel.ID()))).FirstOrCreate(); err != nil {
				return err
			}
		}
		return next(e)
	}
}
