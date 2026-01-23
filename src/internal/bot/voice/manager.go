package voice

import (
	"sync"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

type Manager struct {
	players map[string]*VoicePlayer
	mu      sync.Mutex
}

var manager *Manager

func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			players: make(map[string]*VoicePlayer),
		}
	}
	return manager
}

// 既存の VoicePlayer を返す。なければ nil
func (m *Manager) Get(guildID string) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.players[guildID]
}

// 既存なら返す、なければ新規作成して返す
func (m *Manager) GetOrCreate(guildID string, vc *discordgo.VoiceConnection, ctx *internal.BotContext) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		return p
	}

	p := NewVoicePlayer(guildID, vc, ctx)
	m.players[guildID] = p
	return p
}

// ギルドの VoicePlayer を削除
func (m *Manager) Delete(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		p.Close()
		delete(m.players, guildID)
	}
}
