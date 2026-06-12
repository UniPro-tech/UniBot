package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"slices"
	"sort"
	"unibot/internal/db"
	"unibot/internal/repository"

	"github.com/mmcdole/gofeed"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

	client, err := disgo.New(token,
		//bot.WithDefaultGateway(),
		bot.WithGatewayConfigOpts(
			// Intents
			gateway.WithIntents(
				gateway.IntentsNonPrivileged,
			),
		),
		// Cache
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates),
			cache.WithCaches(cache.FlagChannels),
			cache.WithCaches(cache.FlagMessages),
			cache.WithCaches(cache.FlagRoles),
			cache.WithCaches(cache.FlagMembers),
			cache.WithCaches(cache.FlagGuilds),
		),
		// Event Handler
		bot.WithEventListenerFunc(func(e *events.Ready) {
			Ready(dbConnection, e)
		}),
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
}

func Ready(db *gorm.DB, e *events.Ready) {
	log.Println("Bot is ready 🚀")
	log.Printf("Logged in as: %v#%v", e.User.Username, e.User.Discriminator)

	repo := repository.NewRSSSettingRepository(db)
	rssSubscribeList, err := repo.List()
	if err != nil {
		log.Fatal("An error occured:", err)
	} else {
		for _, rssSetting := range rssSubscribeList {
			// fetch
			url := rssSetting.URL
			fp := gofeed.NewParser()
			feed, err := fp.ParseURL(url)
			if err != nil {
				rssSetting.IsFailed = true
				err := repo.Update(rssSetting)
				if err != nil {
					log.Print("Update Record Failed", err)
					return
				}
			}
			feedTitle := rssSetting.Title
			if feedTitle == nil {
				feedTitle = &feed.Title
			}
			sort.Slice(feed.Items, func(i, j int) bool {
				prev := feed.Items[i]
				next := feed.Items[j]
				if prev.PublishedParsed != nil && next.PublishedParsed != nil {
					prevNano := prev.PublishedParsed.UnixNano()
					nextNano := next.PublishedParsed.UnixNano()
					return prevNano >= nextNano
				}
				return false
			})
			var targetItems []*gofeed.Item
			for _, item := range feed.Items {
				hash := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", item.Title, item.Description)))
				if hash == *rssSetting.LastItemTitleDescriptionHash {
					targetItems = append(targetItems, item)
				}
			}
			slices.Reverse(targetItems)
			channelID := snowflake.MustParse(rssSetting.ChannelID)
			client := e.Client()
			for _, item := range targetItems {
				itemTitle := item.Title
				if itemTitle == "" {
					itemTitle = "(タイトルなし)"
				}
				itemDescription := item.Description
				if itemDescription == "" {
					if item.Content == "" {
						itemDescription = "説明なし"
					} else {
						itemDescription = item.Content
					}
				}
				itemLink := item.Link
				if itemLink == "" {
					itemLink = "リンクなし"
				}
				message := fmt.Sprintf(`# %s に新しい記事が追加されました！
				## %s
				%s
				URL: %s`, *feedTitle, itemTitle, itemDescription, itemLink)
				_, err := client.Rest.CreateMessage(channelID, discord.NewMessageCreate().WithContent(message))
				if err != nil {
					log.Print("Message create error:", err)
				}
			}
		}
	}
}
