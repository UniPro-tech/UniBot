package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"unibot/internal/bot/command"
	"unibot/internal/bot/handler"
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

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is ready 🚀")
	})

	// Commands handler
	dg.AddHandler(handler.InteractionCreate)

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
