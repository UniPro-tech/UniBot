package voice

import (
	"sync"
)

type VoiceManager struct {
	players map[string]*VoicePlayer
	mu      sync.Mutex
}

var manager = &VoiceManager{
	players: map[string]*VoicePlayer{},
}

func GetManager() *VoiceManager {
	return manager
}

func (m *VoiceManager) Get(guildID string) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.players[guildID]
}

func (m *VoiceManager) Set(guildID string, p *VoicePlayer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.players[guildID] = p
}

func (m *VoiceManager) Delete(guildID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.players, guildID)
}
