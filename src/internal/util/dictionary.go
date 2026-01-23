package util

import (
	"log"
	"strings"

	"unibot/internal/model"
	"unibot/internal/repository"

	"gorm.io/gorm"
)

// 辞書を適用する関数
func ApplyDictionary(db *gorm.DB, guildID, content string) string {
	repo := repository.NewTTSDictionaryRepository(db)

	entries, err := repo.ListByGuild(guildID)
	if err != nil {
		log.Println("辞書の取得に失敗しました:", err)
		return content
	}

	return ApplyDictionaryEntries(content, entries)
}

// 辞書エントリを適用する関数（テスト用に分離）
func ApplyDictionaryEntries(content string, entries []*model.TTSDictionary) string {
	for _, entry := range entries {
		if entry.CaseSensitive {
			// 大文字小文字を区別して置換
			content = strings.ReplaceAll(content, entry.Word, entry.Definition)
		} else {
			// 大文字小文字を区別せずに置換
			content = replaceIgnoreCase(content, entry.Word, entry.Definition)
		}
	}

	return content
}

// 大文字小文字を無視して置換する関数
func replaceIgnoreCase(input, old, new string) string {
	lowerInput := strings.ToLower(input)
	lowerOld := strings.ToLower(old)

	var result strings.Builder
	i := 0

	for i < len(input) {
		idx := strings.Index(lowerInput[i:], lowerOld)
		if idx == -1 {
			result.WriteString(input[i:])
			break
		}

		result.WriteString(input[i : i+idx])
		result.WriteString(new)
		i += idx + len(old)
	}

	return result.String()
}
