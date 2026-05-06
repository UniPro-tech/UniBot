package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/godave/golibdave"
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/gorm/logger"

	"unibot/internal"
	"unibot/internal/api/voicevox"
	customHandlers "unibot/internal/bot/handlers/event"
	interaction_handler "unibot/internal/bot/handlers/interaction"
	"unibot/internal/bot/handlers/interaction/command"
	"unibot/internal/db"
)

func main() {
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	dbConnection.Logger = dbConnection.Logger.LogMode(logger.Info)

	ctxData := &internal.BotContext{
		DB:       dbConnection,
		Config:   internal.LoadConfig(),
		VoiceVox: voicevox.New(internal.LoadConfig().VoiceVoxURI, internal.LoadConfig().VoiceVoxAPIKey),
	}

	err = db.SetupDB(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	r := handler.New()

	interaction_handler.RegistHandler(r, ctxData)

	// disgo クライアントの構築
	client, err := disgo.New(token,
		//bot.WithDefaultGateway(),
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
		bot.WithEventListenerFunc(func(e *events.Ready) {
			customHandlers.Ready(ctxData, e)
		}),
		bot.WithEventListenerFunc(func(e *events.MessageCreate) {
			customHandlers.MessageCreate(ctxData, e)
		}),
		bot.WithEventListenerFunc(func(e *events.GuildVoiceStateUpdate) {
			customHandlers.VoiceStateUpdate(ctxData, e)
		}),
		// Cache
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
			cache.WithCaches(cache.FlagChannels),
			cache.WithCaches(cache.FlagMessages),
			cache.WithCaches(cache.FlagRoles),
			cache.WithCaches(cache.FlagMembers),
		),
		// Handler
		bot.WithEventListeners(r),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
	}

	defer client.Close(context.TODO())

	// 接続開始
	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}

	log.Println("Bot is running...")

	// Slash Command の登録
	var generalCommands []discord.ApplicationCommandCreate
	for _, cmd := range command.GeneralCommands {
		generalCommands = append(generalCommands, cmd)
	}

	var adminCommands []discord.ApplicationCommandCreate
	for _, cmd := range command.AdminCommands {
		adminCommands = append(adminCommands, cmd)
	}

	if _, err := client.Rest.SetGlobalCommands(client.ApplicationID, generalCommands); err != nil {
	}
	if _, err := client.Rest.SetGuildCommands(client.ApplicationID, snowflake.MustParse(ctxData.Config.AdminGuildID), adminCommands); err != nil {
		log.Fatal(err)
	}

	// 終了待機
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
