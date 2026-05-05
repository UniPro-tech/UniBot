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
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/gorm/logger"

	"unibot/internal"
	"unibot/internal/api/voicevox"
	"unibot/internal/bot/command"
	"unibot/internal/bot/handler"
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

	// disgo クライアントの構築
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds|gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
		),
		bot.WithEventListenerFunc[events.Ready](func(e events.Ready) {
			handler.Ready(ctxData, &e)
		}),
		bot.WithEventListenerFunc[events.MessageCreate](func(e events.MessageCreate) {
			handler.MessageCreate(ctxData, &e)
		}),
		bot.WithEventListenerFunc[events.ApplicationCommandInteractionCreate](func(e events.ApplicationCommandInteractionCreate) {
			handler.InteractionCreate(ctxData)(&e)
		}),
		bot.WithEventListenerFunc[events.GuildVoiceStateUpdate](func(e events.GuildVoiceStateUpdate) {
			handler.VoiceStateUpdate(ctxData, &e)
		}),
	)
	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
	}

	// 接続開始
	if err = client.OpenGateway(context.TODO()); err != nil {
		log.Fatal("error while connecting to gateway: ", err)
	}
	defer client.Close(context.TODO())

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
		log.Fatalf("error while registering commands: %v", err)
	}
	if _, err := client.Rest.SetGuildCommands(client.ApplicationID, snowflake.MustParse(ctxData.Config.AdminGuildID), adminCommands); err != nil {
		log.Fatalf("error while registering commands: %v", err)
	}

	// 終了待機
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
