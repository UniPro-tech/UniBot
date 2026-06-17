package interaction_handler

import (
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

func IsOwner(member discord.Member) bool {
	config := internal.LoadConfig()
	adminRoleID := config.AdminRoleID
	return slices.Contains(member.RoleIDs, snowflake.MustParse(adminRoleID))
}

func AdminOnlyMiddleware(ctx *internal.BotContext) func(next handler.Handler) handler.Handler {
	return func(next handler.Handler) handler.Handler {
		return func(e *handler.InteractionEvent) error {
			config := ctx.Config
			if !IsOwner(e.Member().Member) {
				errorEmbed := discord.Embed{
					Title:       "権限エラー",
					Description: "権限がありません。",
					Color:       config.Colors.Error,
					Footer: &discord.EmbedFooter{
						Text:    "Requested by " + *e.Member().Nick,
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
				e.DeferCreateMessage(ephemeral)
				return next(e)
			}
		} else {
			return func(e *handler.InteractionEvent) error {
				e.DeferUpdateMessage()
				return next(e)
			}
		}
	}
}
