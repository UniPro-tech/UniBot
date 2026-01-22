package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

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

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	ctx := &internal.BotContext{
		DB:       dbConnection,
		Config:   internal.LoadConfig(),
		VoiceVox: voicevox.New(internal.LoadConfig().VoiceVoxURI, internal.LoadConfig().VoiceVoxAPIKey),
	}

	err = db.SetupDB(dbConnection)
	if err != nil {
		log.Fatal(err)
	}

	// Register handlers
	dg.AddHandler(handler.Ready(ctx))
	dg.AddHandler(handler.MessageCreate(ctx))
	dg.AddHandler(handler.InteractionCreate(ctx))
	dg.AddHandler(handler.VoiceStateUpdate(ctx))

	// Start the bot
	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dg.Close()

	log.Println("Bot is running...")

	// Register commands
	appID := dg.State.User.ID

	for _, cmd := range command.Commands {
		_, err := dg.ApplicationCommandCreate(appID, "", cmd)
		if err != nil {
			log.Fatalf("cannot create command %s: %v", cmd.Name, err)
		}
	}

	// Exit handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}
