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
	"unibot/internal/util"

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
		return
	}
	for _, rssSetting := range rssSubscribeList {
		url := rssSetting.URL
		feed, err := util.FetchFeed(url)
		if err != nil {
			rssSetting.IsFailed = true
			if err := repo.Update(rssSetting); err != nil {
				log.Print("Update Record Failed", err)
			}
			continue
		}

		feedTitle := rssSetting.Title
		if feedTitle == nil {
			feedTitle = &feed.Title
		}

		// 新しい日時がindex:0
		sort.Slice(feed.Items, func(i, j int) bool {
			prev := feed.Items[i]
			next := feed.Items[j]
			if prev.PublishedParsed != nil && next.PublishedParsed != nil {
				return prev.PublishedParsed.UnixNano() >= next.PublishedParsed.UnixNano()
			}
			return false
		})

		// 保存済みハッシュより新しい記事を収集する
		// - LastItemTitleDescriptionHash が nil（初回）なら全件対象
		// - 一致するハッシュが見つかった時点で break（それ以降は既読）
		var targetItems []*gofeed.Item
		for _, item := range feed.Items {
			hash := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", item.Title, item.Description)))
			if rssSetting.LastItemTitleDescriptionHash != nil && hash == *rssSetting.LastItemTitleDescriptionHash {
				break // ここから先は既読
			}
			targetItems = append(targetItems, item)
		}

		if len(targetItems) == 0 {
			continue
		}

		// 古い→新しい順に送信
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
				if item.Content != "" {
					itemDescription = item.Content
				} else {
					itemDescription = "(説明なし)"
				}
			}
			itemLink := item.Link
			if itemLink == "" {
				itemLink = "リンクなし"
			}
			message := fmt.Sprintf(
				"# %s に新しい記事が追加されました！\n## %s\n%s\nURL: %s",
				*feedTitle, itemTitle, itemDescription, itemLink,
			)
			_, err := client.Rest.CreateMessage(channelID, discord.NewMessageCreate().WithContent(message))
			if err != nil {
				log.Print("Message create error:", err)
			}
		}

		// 送信完了後に最新ハッシュを保存（feed.Items[0] = 最新記事）
		newestHash := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", feed.Items[0].Title, feed.Items[0].Description)))
		rssSetting.LastItemTitleDescriptionHash = &newestHash
		rssSetting.IsFailed = false
		if err := repo.Update(rssSetting); err != nil {
			log.Print("Update Record Failed", err)
		}
	}
}
