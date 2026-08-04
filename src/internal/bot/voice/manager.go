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

// GetManager は sync.Once を用いて安全にシングルトンインスタンスを返却する
// これにより、複数の goroutine から同時にアクセスされても、
// Manager のインスタンスが一度だけ生成されることが保証される
func GetManager() *Manager {
	managerOnce.Do(func() {
		managerInstance = &Manager{
			players: make(map[int64]*VoicePlayer),
		}
	})
	return managerInstance
}

// GetPlayer は、指定されたギルドIDに対応する VoicePlayer を取得する。
// もし存在しない場合は nil を返す。
// この関数はスレッドセーフであり、複数の goroutine から同時に呼び出されても安全に動作する。
func (m *Manager) GetPlayer(guildID int64) *VoicePlayer {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.players[guildID]
}

// GetOrCreatePlayer は、指定されたギルドIDに対応する VoicePlayer を取得する。
// もし存在しない場合は、新しい VoicePlayer を作成して返す。
// この関数はスレッドセーフであり、複数の goroutine から同時に呼び出されても安全に動作する。
func (m *Manager) GetOrCreatePlayer(
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

// Delete は、指定されたギルドIDに対応する VoicePlayer を削除する。
// もし存在しない場合は何も行わない。
// この関数はスレッドセーフであり、複数の goroutine から同時に呼び出されても安全に動作する。
func (m *Manager) DeletePlayer(guildID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p, ok := m.players[guildID]; ok {
		p.Close()
		delete(m.players, guildID)
	}
}
