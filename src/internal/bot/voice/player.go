package voice

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
	"unibot/internal"
	"unibot/internal/model"

	"github.com/disgoorg/disgo/voice"
	"github.com/hraban/opus"
)

type QueueItem struct {
	Text    string
	Setting model.TTSPersonalSetting
}

type VoicePlayer struct {
	GuildID    string
	ChannelID  string
	VC         voice.Conn
	playing    bool
	initFrames int
	vcMu       sync.RWMutex
	TextQueue  chan QueueItem
	Stop       chan struct{}
	Skip       chan struct{}

	// 音声データ供給用の状態管理
	currentOpusChan chan []byte
	currentCtx      context.Context
	currentCancel   context.CancelFunc
}

func NewVoicePlayer(guildID string, channelID string, vc voice.Conn, ctx *internal.BotContext) *VoicePlayer {
	p := &VoicePlayer{
		GuildID:         guildID,
		ChannelID:       channelID,
		VC:              vc,
		TextQueue:       make(chan QueueItem, 50),
		Stop:            make(chan struct{}),
		Skip:            make(chan struct{}),
		currentOpusChan: make(chan []byte, 10),
		playing:         false,
	}

	go p.worker(ctx)
	return p
}

func (p *VoicePlayer) ProvideOpusFrame() ([]byte, error) {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	if !p.playing {
		return nil, nil
	}

	// パニック回避策:
	// ProvideOpusFrameが呼ばれ始めてから、内部のUDPが安定するまで
	// 最初の数フレーム（約100ms分）をあえて nil で返す
	if p.initFrames < 5 {
		p.initFrames++
		return nil, nil
	}

	select {
	case frame := <-p.currentOpusChan:
		return frame, nil
	default:
		return nil, nil
	}
}

// 以下の2つはインターフェースを満たすために必要

func (p *VoicePlayer) Close() {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	select {
	case <-p.Stop:
		// すでに閉じている
	default:
		close(p.Stop)
		close(p.TextQueue) // これを閉じると worker の for range が即座に終了する
	}
}

// これが true の時だけ Disgo は ProvideOpusFrame を呼びに来る
func (p *VoicePlayer) CanProvide() bool {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	// 再生フラグが立っていない時は供給しない
	if !p.playing {
		return false
	}

	// 接続状態をチェック
	vc := p.GetVC()
	// Disgoの内部で接続が完全に確立（UDP/DAVE完了）しているかを確認
	// 接続がnil、または有効でない場合はfalseを返してDisgoの送信ループを待機させる
	return vc != nil && vc.ChannelID() != nil
}

func (p *VoicePlayer) SetVC(vc voice.Conn) {
	p.vcMu.Lock()
	p.VC = vc
	if vc != nil {
		vc.SetOpusFrameProvider(p)
	}
	p.vcMu.Unlock()
}

func (p *VoicePlayer) worker(ctx *internal.BotContext) {
	for item := range p.TextQueue {
		p.currentCtx, p.currentCancel = context.WithCancel(context.Background())

		// 1. 合成
		audio, err := ctx.VoiceVox.Synthesize(p.currentCtx, item.Text, item.Setting.SpeakerID, float64(item.Setting.SpeakerSpeed)/100.0)
		if err != nil {
			p.currentCancel()
			continue
		}

		vc := p.GetVC()
		if vc == nil {
			p.currentCancel()
			continue
		}

		// 2. 接続の生存確認
		if vc.ChannelID() == nil {
			p.currentCancel()
			continue
		}

		// 3. Providerをセットする直前に、Speak状態にして少し待つ
		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagMicrophone)
		time.Sleep(200 * time.Millisecond) // UDPソケットが安定するための「祈り」の時間

		p.vcMu.Lock()
		p.initFrames = 0 // カウンタをリセット
		p.playing = true
		p.vcMu.Unlock()

		// ここで登録
		vc.SetOpusFrameProvider(p)

		// 4. 再生
		p.streamAudio(p.currentCtx, audio)

		// 5. 終了
		p.vcMu.Lock()
		p.playing = false
		p.vcMu.Unlock()

		vc.SetOpusFrameProvider(nil)
		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)
		p.currentCancel()
		select {
		case <-p.Stop:
			return
		default:
		}
	}
}

func (p *VoicePlayer) streamAudio(ctx context.Context, wav []byte) {
	tmp, _ := os.CreateTemp("", "tts-*.wav")
	defer os.Remove(tmp.Name())
	tmp.Write(wav)
	tmp.Close()

	cmd := exec.Command("ffmpeg", "-loglevel", "quiet", "-i", tmp.Name(), "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return
	}
	defer cmd.Process.Kill()

	enc, _ := opus.NewEncoder(48000, 2, opus.AppAudio)
	pcm := make([]int16, 960*2)
	byteBuf := make([]byte, len(pcm)*2)

	for {
		_, err := io.ReadFull(stdout, byteBuf)
		if err != nil {
			break
		}

		for i := 0; i < len(pcm); i++ {
			pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
		}

		opusBuf := make([]byte, 4000)
		n, _ := enc.Encode(pcm, opusBuf)

		select {
		case <-ctx.Done():
			return
		case <-p.Skip:
			return
		case p.currentOpusChan <- opusBuf[:n]:
			// Providerが消費するまでここで待機（実質的にタイミング制御になる）
		}
	}
	// ffmpeg の読み取りが終わった後、5フレーム分(約100ms)の無音を送って、
	// 最後の一言がバッファで切られないように押し出す
	silenceFrame := []byte{0xF8, 0xFF, 0xFE} // Opus の標準的な無音フレーム
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return
		case p.currentOpusChan <- silenceFrame:
		}
	}
}

func (p *VoicePlayer) EnqueueText(item QueueItem) {
	select {
	case p.TextQueue <- item:
	default:
	}
}

func (p *VoicePlayer) SkipCurrent() {
	select {
	case p.Skip <- struct{}{}:
	default:
	}
}

func (p *VoicePlayer) GetVC() voice.Conn {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()
	return p.VC
}
