package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"unibot/internal/model"
	"unibot/internal/repository"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&model.TTSDictionary{})
	require.NoError(t, err)

	return db
}

// TypeScript版と同じ動作をするかテスト

// テスト1: 基本的な置換
func TestApplyDictionary_BasicReplacement(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "テスト",
		Definition:    "てすと",
		CaseSensitive: false,
	})

	input := "これはテストです"
	expected := "これはてすとです"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト2: 大文字小文字を区別しない置換 (TypeScript版: caseSensitive: false -> "gi" フラグ)
func TestApplyDictionary_CaseInsensitive(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "Hello",
		Definition:    "こんにちは",
		CaseSensitive: false,
	})

	// TypeScript版: new RegExp("Hello", "gi") で HELLO, hello, HeLLo 全部マッチ
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello world", "こんにちは world"},
		{"HELLO world", "こんにちは world"},
		{"hello world", "こんにちは world"},
		{"HeLLo world", "こんにちは world"},
	}

	for _, tc := range testCases {
		result := ApplyDictionary(db, "guild123", tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}

// テスト3: 大文字小文字を区別する置換 (TypeScript版: caseSensitive: true -> "g" フラグ)
func TestApplyDictionary_CaseSensitive(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "Hello",
		Definition:    "こんにちは",
		CaseSensitive: true,
	})

	// TypeScript版: new RegExp("Hello", "g") で Hello のみマッチ
	testCases := []struct {
		input    string
		expected string
	}{
		{"Hello world", "こんにちは world"},
		{"HELLO world", "HELLO world"},       // マッチしない
		{"hello world", "hello world"},       // マッチしない
		{"HeLLo world", "HeLLo world"},       // マッチしない
	}

	for _, tc := range testCases {
		result := ApplyDictionary(db, "guild123", tc.input)
		assert.Equal(t, tc.expected, result, "Input: %s", tc.input)
	}
}

// テスト4: 複数の辞書エントリ
func TestApplyDictionary_MultipleEntries(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "www",
		Definition:    "わらわらわら",
		CaseSensitive: false,
	})
	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "lol",
		Definition:    "笑",
		CaseSensitive: false,
	})

	input := "www それな lol"
	expected := "わらわらわら それな 笑"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト5: 同じ単語が複数回出現
func TestApplyDictionary_MultipleOccurrences(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "草",
		Definition:    "くさ",
		CaseSensitive: false,
	})

	// TypeScript版: "g" フラグで全置換
	input := "草草草"
	expected := "くさくさくさ"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト6: 辞書が空の場合
func TestApplyDictionary_EmptyDictionary(t *testing.T) {
	db := setupTestDB(t)

	input := "テキストそのまま"
	expected := "テキストそのまま"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト7: 別のギルドの辞書は適用されない
func TestApplyDictionary_DifferentGuild(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild999",
		UserID:        "user1",
		Word:          "テスト",
		Definition:    "てすと",
		CaseSensitive: false,
	})

	input := "これはテストです"
	expected := "これはテストです" // 置換されない
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト8: 日本語の大文字小文字（実際は関係ないがテスト）
func TestApplyDictionary_JapaneseText(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "ありがとう",
		Definition:    "あざす",
		CaseSensitive: false,
	})

	input := "ありがとうございます"
	expected := "あざすございます"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト9: 特殊文字を含む単語
func TestApplyDictionary_SpecialCharacters(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "(笑)",
		Definition:    "わら",
		CaseSensitive: false,
	})

	input := "面白い(笑)"
	expected := "面白いわら"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト10: 置換後のテキストが別の辞書エントリにマッチするケース
func TestApplyDictionary_ChainedReplacement(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewTTSDictionaryRepository(db)

	// TypeScript版は順番に適用される
	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "A",
		Definition:    "B",
		CaseSensitive: true,
	})
	repo.Create(&model.TTSDictionary{
		GuildID:       "guild123",
		UserID:        "user1",
		Word:          "B",
		Definition:    "C",
		CaseSensitive: true,
	})

	input := "A"
	// 1回目: A -> B, 2回目: B -> C
	// TypeScript版と同じく、順番に適用されるので A -> B -> C
	expected := "C"
	result := ApplyDictionary(db, "guild123", input)

	assert.Equal(t, expected, result)
}

// テスト11: replaceIgnoreCase 単体テスト
func TestReplaceIgnoreCase(t *testing.T) {
	testCases := []struct {
		input    string
		old      string
		new      string
		expected string
	}{
		{"Hello World", "hello", "Hi", "Hi World"},
		{"HELLO World", "hello", "Hi", "Hi World"},
		{"hello World", "hello", "Hi", "Hi World"},
		{"HeLLo World", "hello", "Hi", "Hi World"},
		{"Hello hello HELLO", "hello", "Hi", "Hi Hi Hi"},
		{"No match here", "xyz", "abc", "No match here"},
		{"", "test", "result", ""},
	}

	for _, tc := range testCases {
		result := replaceIgnoreCase(tc.input, tc.old, tc.new)
		assert.Equal(t, tc.expected, result, "Input: %s, Old: %s", tc.input, tc.old)
	}
}
