package voice

import (
	"sync"
	"unibot/internal"

	"github.com/disgoorg/disgo/voice"
)

type Manager struct {
	players map[string]*VoicePlayer
	mu      sync.Mutex
}

var (
	manager *Manager
	once    sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		manager = &Manager{
			players: make(map[string]*VoicePlayer),
		}
	})
	return manager
}

func (m *Manager) Get(guildID string) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.players[guildID]
}

func (m *Manager) GetOrCreate(
	guildID string,
	channelID string,
	vc voice.Conn,
	ctx *internal.BotContext,
) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		// すでにプレイヤーが存在する場合、新しい接続があれば更新
		if vc != nil {
			p.SetVC(vc)
			p.ChannelID = channelID
		}
		return p
	}

	p := NewVoicePlayer(guildID, channelID, vc, ctx)
	m.players[guildID] = p
	return p
}

func (m *Manager) Delete(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		p.Close()
		delete(m.players, guildID)
	}
}
