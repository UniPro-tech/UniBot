package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"slices"
	"sort"
	"unibot/internal/db"
	"unibot/internal/repository"
	"unibot/internal/util"

	"github.com/mmcdole/gofeed"

	"github.com/disgoorg/disgo/webhook"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	dbConnection.Logger = dbConnection.Logger.LogMode(logger.Info)

	log.Println("Bot is running...")

	Ready(dbConnection)
}

func Ready(db *gorm.DB) {
	log.Println("Bot is ready 🚀")

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
		client, err := webhook.NewWithURL(rssSetting.WebhookURL)
		if err != nil {
			log.Print("Message create error:", err)
			rssSetting.IsFailed = true
			if err := repo.Update(rssSetting); err != nil {
				log.Print("Update Record Failed", err)
			}
			continue
		}
		for _, item := range targetItems {
			itemTitle := item.Title
			if itemTitle == "" {
				itemTitle = "(タイトルなし)"
			}
			itemLink := item.Link
			if itemLink == "" {
				itemLink = "リンクなし"
			}
			var itemAuthor string
			if len(item.Authors) > 0 {
				for index, author := range item.Authors {
					var separate string
					if index != 0 {
						separate = ","
					}
					itemAuthor = itemAuthor + separate + author.Name
				}
			}
			itemPublished := item.PublishedParsed

			message := fmt.Sprintf(
				"# %s に新しい記事が追加されました！\n## %s\n-# by %s at <t:%d:S>\nURL: %s",
				*feedTitle, itemTitle, itemAuthor, itemPublished.UTC().Unix(), itemLink,
			)
			_, err := client.CreateContent(message)
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
