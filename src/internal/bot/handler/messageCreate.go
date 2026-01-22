package handler

import (
	"log"
	"unibot/internal/db"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, r *discordgo.MessageCreate) {
	// Ignore messages from the bot itself
	if r.Author.ID == s.State.User.ID {
		return
	}

	// If the guild is nil, it's a DM; ignore it
	if r.GuildID == "" {
		return
	}

	// If bot is connected to a voice channel, TTS the message
	dbConnection, err := db.NewDB()
	if err != nil {
		log.Print(err)
	}

	repo := repository.NewTTSConnectionRepository(dbConnection)
	repo.GetByChannelID(r.ChannelID)
}
