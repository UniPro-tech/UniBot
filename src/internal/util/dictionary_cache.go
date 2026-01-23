package util

import (
	"reflect"
	"sync"
	"time"

	"unibot/internal/model"
	"unibot/internal/repository"

	"gorm.io/gorm"
)

// DictionaryCache はギルドごとの辞書エントリをキャッシュします
type DictionaryCache struct {
	mu      sync.RWMutex
	entries map[cacheKey]*cacheEntry
	inFlight map[cacheKey]*inFlight
	ttl     time.Duration
}

type cacheEntry struct {
	entries   []*model.TTSDictionary
	expiresAt time.Time
}

type cacheKey struct {
	dbPtr   uintptr
	guildID string
}

type inFlight struct {
	wg      sync.WaitGroup
	entries []*model.TTSDictionary
	err     error
}

// DefaultCacheTTL はデフォルトのキャッシュ有効期限
const DefaultCacheTTL = 5 * time.Minute

var (
	globalCache     *DictionaryCache
	globalCacheOnce sync.Once
)

// GetDictionaryCache はグローバルなキャッシュインスタンスを返します
func GetDictionaryCache() *DictionaryCache {
	globalCacheOnce.Do(func() {
		globalCache = NewDictionaryCache(DefaultCacheTTL)
	})
	return globalCache
}

// NewDictionaryCache は新しいキャッシュインスタンスを作成します
func NewDictionaryCache(ttl time.Duration) *DictionaryCache {
	return &DictionaryCache{
		entries:  make(map[cacheKey]*cacheEntry),
		inFlight: make(map[cacheKey]*inFlight),
		ttl:      ttl,
	}
}

// Get はキャッシュから辞書エントリを取得します
// キャッシュミスまたは期限切れの場合はDBから取得し、キャッシュを更新します
func (c *DictionaryCache) Get(db *gorm.DB, guildID string) ([]*model.TTSDictionary, error) {
	key := cacheKeyFor(db, guildID)
	now := time.Now()

	c.mu.RLock()
	entry, exists := c.entries[key]
	if exists && now.Before(entry.expiresAt) {
		c.mu.RUnlock()
		return entry.entries, nil
	}
	if flight, ok := c.inFlight[key]; ok {
		c.mu.RUnlock()
		flight.wg.Wait()
		return flight.entries, flight.err
	}
	c.mu.RUnlock()

	c.mu.Lock()
	entry, exists = c.entries[key]
	if exists && time.Now().Before(entry.expiresAt) {
		c.mu.Unlock()
		return entry.entries, nil
	}
	if flight, ok := c.inFlight[key]; ok {
		c.mu.Unlock()
		flight.wg.Wait()
		return flight.entries, flight.err
	}

	flight := &inFlight{}
	flight.wg.Add(1)
	c.inFlight[key] = flight
	c.mu.Unlock()

	entries, err := c.fetchFromDB(db, guildID)

	c.mu.Lock()
	if current, ok := c.inFlight[key]; ok && current == flight {
		if err == nil {
			c.entries[key] = &cacheEntry{
				entries:   entries,
				expiresAt: time.Now().Add(c.ttl),
			}
		}
		delete(c.inFlight, key)
	}
	flight.entries = entries
	flight.err = err
	flight.wg.Done()
	c.mu.Unlock()

	return entries, err
}

// fetchFromDB はDBから辞書エントリを取得します
func (c *DictionaryCache) fetchFromDB(db *gorm.DB, guildID string) ([]*model.TTSDictionary, error) {
	repo := repository.NewTTSDictionaryRepository(db)
	entries, err := repo.ListByGuild(guildID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

// Invalidate は指定したギルドのキャッシュを無効化します
func (c *DictionaryCache) Invalidate(guildID string) {
	c.mu.Lock()
	for key := range c.entries {
		if key.guildID == guildID {
			delete(c.entries, key)
		}
	}
	for key := range c.inFlight {
		if key.guildID == guildID {
			delete(c.inFlight, key)
		}
	}
	c.mu.Unlock()
}

// InvalidateAll はすべてのキャッシュを無効化します
func (c *DictionaryCache) InvalidateAll() {
	c.mu.Lock()
	c.entries = make(map[cacheKey]*cacheEntry)
	c.inFlight = make(map[cacheKey]*inFlight)
	c.mu.Unlock()
}

func cacheKeyFor(db *gorm.DB, guildID string) cacheKey {
	return cacheKey{
		dbPtr:   reflect.ValueOf(db).Pointer(),
		guildID: guildID,
	}
}
