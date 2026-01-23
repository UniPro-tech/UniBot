package util

import (
	"log"
	"strings"

	"unibot/internal/model"

	"gorm.io/gorm"
)

// 辞書を適用する関数（キャッシュを使用）
func ApplyDictionary(db *gorm.DB, guildID, content string) string {
	cache := GetDictionaryCache()

	entries, err := cache.Get(db, guildID)
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
	// 空文字列の場合は無限ループを防ぐためそのまま返す
	if old == "" {
		return input
	}

	lowerInput := strings.ToLower(input)
	lowerOld := strings.ToLower(old)

	// マルチバイト文字対応: 小文字化後のバイト長が異なる場合があるため
	// rune単位で処理する
	inputRunes := []rune(input)
	lowerInputRunes := []rune(lowerInput)
	lowerOldRunes := []rune(lowerOld)
	oldRuneLen := len([]rune(old))

	var result strings.Builder
	i := 0

	for i < len(inputRunes) {
		idx := indexRunes(lowerInputRunes[i:], lowerOldRunes)
		if idx == -1 {
			result.WriteString(string(inputRunes[i:]))
			break
		}

		result.WriteString(string(inputRunes[i : i+idx]))
		result.WriteString(new)
		i += idx + oldRuneLen
	}

	return result.String()
}

// rune配列から部分配列を検索する関数
func indexRunes(s, substr []rune) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

Outer:
	for i := 0; i <= len(s)-len(substr); i++ {
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				continue Outer
			}
		}
		return i
	}
	return -1
}
