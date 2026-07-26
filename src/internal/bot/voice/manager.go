package voice

import (
	"sync"
	"unibot/internal"

	"github.com/disgoorg/disgo/voice"
)

type Manager struct {
	players map[int64]*VoicePlayer
	mu      sync.Mutex
}

var (
	managerInstance *Manager
	managerOnce     sync.Once
)

// GetManager は sync.Once を用いて安全にシングルトンインスタンスを返します
func GetManager() *Manager {
	managerOnce.Do(func() {
		managerInstance = &Manager{
			players: make(map[int64]*VoicePlayer),
		}
	})
	return managerInstance
}

func (m *Manager) Get(guildID int64) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.players[guildID]
}

func (m *Manager) GetOrCreate(
	guildID int64,
	channelID int64,
	vc voice.Conn,
	ctx *internal.BotContext,
) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
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

func (m *Manager) Delete(guildID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		p.Close()
		delete(m.players, guildID)
	}
}
