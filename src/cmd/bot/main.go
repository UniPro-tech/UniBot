package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/godave/golibdave"
	"github.com/disgoorg/snowflake/v2"

	"unibot/internal"
	"unibot/internal/api/voicevox"
	customHandlers "unibot/internal/bot/handlers/event"
	interaction_handler "unibot/internal/bot/handlers/interaction"
	"unibot/internal/bot/handlers/interaction/command"
	"unibot/internal/bot/notify"
	"unibot/internal/db"
	"unibot/internal/logger"
	"unibot/internal/logger/discordsink"
	"unibot/internal/logger/gormslog"
	"unibot/internal/query"
	"unibot/internal/util"
)

// 終了時にログを吐き出すための猶予。
const shutdownTimeout = 10 * time.Second

func main() {
	envErr := util.LoadProjectEnv()

	// Discord への通知シンクはこの時点では送信先を持たない。
	// disgo クライアント構築後に Attach する。
	sink := discordsink.New(discordsink.ConfigFromEnv())

	logOpts := logger.ConfigFromEnv()
	logOpts.Version = internal.Version
	logOpts.Commit = internal.GitCommit
	logOpts.Branch = internal.GitBranch
	logOpts.Extra = []slog.Handler{sink}
	logger.Init(logOpts)

	ctx := context.Background()

	if envErr != nil {
		// コンテナ上では .env が存在しないのが正常なので Debug に留める。
		slog.DebugContext(ctx, "no .env file loaded", slog.Any("err", envErr))
	}

	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		fatal(ctx, "DISCORD_TOKEN is not set")
	}

	dbConnection, err := db.NewDB(gormslog.New(gormslog.ConfigFromEnv()))
	if err != nil {
		fatal(ctx, "failed to connect to database", slog.Any("err", err))
	}

	config := internal.LoadConfig()
	ctxData := &internal.BotContext{
		DB:       dbConnection,
		Config:   config,
		VoiceVox: voicevox.New(config.VoiceVoxURI, config.VoiceVoxAPIKey),
	}

	query.SetDefault(dbConnection)

	notify.Init(notify.ConfigFromEnv())

	r := handler.New()
	// インタラクションごとに trace_id を発行する。
	r.DefaultContext(func() context.Context {
		traced, _ := logger.WithNewTrace(context.Background())
		return traced
	})
	// グローバルミドルウェアを経由しない経路のための保険。
	r.Error(interaction_handler.ErrorSink)
	// 最も外側で動かすことで、CreateMasterRecordMiddleware の失敗も捕捉する。
	r.Use(interaction_handler.ObservabilityMiddleware(ctxData))
	r.Use(interaction_handler.CreateMasterRecordMiddleware)

	interaction_handler.RegistHandler(r, ctxData)

	// disgo クライアントの構築
	client, err := disgo.New(token,
		//bot.WithDefaultGateway(),
		// Logger
		// disgo には Discord シンクを持たないロガーを渡す。
		// disgo のレートリミッタは Warn を出すため、シンク付きロガーだと
		// 「通知の送信で 429 → Warn → また通知」という自己増殖ループになる。
		bot.WithLogger(logger.Plain("disgo")),
		bot.WithGatewayConfigOpts(
			// Intents
			gateway.WithIntents(
				gateway.IntentsNonPrivileged,
				gateway.IntentMessageContent,
			),
		),
		// DAVE
		bot.WithVoiceManagerConfigOpts(
			voice.WithDaveSessionCreateFunc(golibdave.NewSession),
		),
		// Event Listener
		bot.WithEventListenerFunc(customHandlers.Wrap("ready", ctxData, customHandlers.Ready)),
		bot.WithEventListenerFunc(customHandlers.Wrap("message_create", ctxData, customHandlers.MessageCreate)),
		bot.WithEventListenerFunc(customHandlers.Wrap("voice_state_update", ctxData, customHandlers.VoiceStateUpdate)),
		// Cache
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
			cache.WithCaches(cache.FlagChannels),
			cache.WithCaches(cache.FlagMessages),
			cache.WithCaches(cache.FlagRoles),
			cache.WithCaches(cache.FlagMembers),
			cache.WithCaches(cache.FlagGuilds),
		),
		// Handler
		bot.WithEventListeners(r),
	)
	if err != nil {
		fatal(ctx, "failed to build disgo client", slog.Any("err", err))
	}

	sink.Attach(client.Rest)

	// 接続開始
	if err = client.OpenGateway(ctx); err != nil {
		fatal(ctx, "failed to open gateway", slog.Any("err", err))
	}

	logger.Notice(ctx, "bot started",
		slog.String("version", ctxData.Config.BotVersion),
		slog.String("commit", internal.GitCommit),
		slog.String("branch", internal.GitBranch),
	)

	registerCommands(ctx, client, ctxData)

	// 終了待機
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	sig := <-stop

	shutdown(client, ctxData, sig)
}

// registerCommands はスラッシュコマンドを Discord に同期する。
// 失敗しても起動は継続する（一般コマンドが使える限りプロセスを落とす必要はない）。
func registerCommands(ctx context.Context, client *bot.Client, ctxData *internal.BotContext) {
	var generalCommands []discord.ApplicationCommandCreate
	for _, cmd := range command.GeneralCommands {
		generalCommands = append(generalCommands, cmd)
	}

	var adminCommands []discord.ApplicationCommandCreate
	for _, cmd := range command.AdminCommands {
		adminCommands = append(adminCommands, cmd)
	}

	if _, err := client.Rest.SetGlobalCommands(client.ApplicationID, generalCommands); err != nil {
		slog.ErrorContext(ctx, "failed to register global commands", slog.Any("err", err))
	} else {
		slog.InfoContext(ctx, "registered global commands", slog.Int("count", len(generalCommands)))
	}

	adminGuildID, err := snowflake.Parse(ctxData.Config.AdminGuildID)
	if err != nil {
		slog.ErrorContext(ctx, "invalid admin guild id, admin commands were not registered",
			slog.String("value", ctxData.Config.AdminGuildID), slog.Any("err", err))
		return
	}

	if _, err := client.Rest.SetGuildCommands(client.ApplicationID, adminGuildID, adminCommands); err != nil {
		slog.ErrorContext(ctx, "failed to register admin commands", slog.Any("err", err))
		return
	}
	slog.InfoContext(ctx, "registered admin commands", slog.Int("count", len(adminCommands)))
}

func shutdown(client *bot.Client, ctxData *internal.BotContext, sig os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	ctx, _ = logger.WithTrace(ctx)

	logger.Notice(ctx, "bot is shutting down", slog.String("signal", sig.String()))

	notify.Shutdown(ctx, client, ctxData.Config, sig.String())

	// client.Close は client.Rest も閉じるため、シンクのフラッシュを先に行う。
	// 順序を逆にすると、キューに残ったログを送信できない。
	if err := logger.Shutdown(ctx); err != nil {
		// シンクが閉じられないだけなので stdout に残して終了する。
		logger.Plain("shutdown").Warn("failed to flush log sink", slog.Any("err", err))
	}

	client.Close(ctx)
}

// fatal はエラーを記録し、Discord シンクを吐き出してからプロセスを終了する。
// os.Exit は defer を実行しないため、フラッシュは明示的に行う必要がある。
func fatal(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)

	flushCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	_ = logger.Shutdown(flushCtx)
	cancel()

	os.Exit(1)
}
