package general

import (
	"log"
	"time"
	"unibot/internal"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func LoadPingCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "ping",
		Description: "スピードテストを行います",
	}
}

func Ping(ctx *internal.BotContext, e *events.ApplicationCommandInteractionCreate) {
	config := ctx.Config

	// Get ws websocketLatency
	websocketLatency := e.Client().Gateway.Latency()

	username := ""
	if e.Member().Nick != nil {
		username = *e.Member().Nick
	} else if e.User().GlobalName != nil {
		username = *e.User().GlobalName
	} else {
		username = e.User().Username
	}

	// Respond to the interaction
	responseEmbed := discord.Embed{
		Title:       "Pong 🏓",
		Description: "スピードテストの結果です",
		Color:       config.Colors.Primary,
		Fields: []discord.EmbedField{
			{
				Name:  "WebSocket Latency",
				Value: websocketLatency.String(),
			},
		},
		Author: &discord.EmbedAuthor{
			IconURL: *e.Member().Avatar,
			Name:    *e.Member().Nick,
		},
		Footer: &discord.EmbedFooter{
			Text:    "Requested by " + username,
			IconURL: *e.Member().Avatar,
		},
		Timestamp: func() *time.Time {
			t := time.Now()
			return &t
		}(),
	}
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			embed := discord.NewEmbed().
				WithTitle("Embed Title").
				WithDescription("This is a description").
				WithColor(0x5865F2)
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
				discord.NewMessageCreate().WithEmbeds(embed))
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)

	_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(),
		discord.NewMessageCreate().WithEmbeds(responseEmbed))
	if err != nil {
		embed := discord.Embed{
			Title:       "エラー",
			Description: "スピードテストに失敗しました。",
			Color:       config.Colors.Error,
			Footer: &discord.EmbedFooter{
				Text:    "Requested by " + username,
				IconURL: *e.Member().Avatar,
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}
		e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEmbeds(embed).WithEphemeral(true))
		return
	}
}
