package voice

import (
	"context"
	"encoding/binary"
	"io"
	"os/exec"
	"sync"
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
	// エンコーダ
	encoder *opus.Encoder
}

func NewVoicePlayer(guildID string, channelID string, vc voice.Conn, ctx *internal.BotContext) *VoicePlayer {
	// エンコーダを作成
	enc, _ := opus.NewEncoder(48000, 2, opus.AppAudio)

	p := &VoicePlayer{
		GuildID:         guildID,
		ChannelID:       channelID,
		VC:              vc,
		TextQueue:       make(chan QueueItem, 50),
		Stop:            make(chan struct{}),
		Skip:            make(chan struct{}),
		currentOpusChan: make(chan []byte, 10),
		playing:         false,
		encoder:         enc,
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
	p.vcMu.RLock() // Lock() ではなく RLock() に変更
	defer p.vcMu.RUnlock()

	// 再生フラグが立っていない時は供給しない
	if !p.playing {
		return false
	}

	vc := p.VC
	if vc == nil || vc.Gateway().Status() != voice.StatusReady {
		return false
	}

	return true
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
		// コンテキストの生成と代入をロックで保護
		cCtx, cCancel := context.WithCancel(context.Background())
		p.vcMu.Lock()
		p.currentCtx = cCtx
		p.currentCancel = cCancel
		p.vcMu.Unlock()

		// 1. 合成 (cCtx を渡すことで、合成中のスキップにも対応可能)
		audio, err := ctx.VoiceVox.Synthesize(cCtx, item.Text, item.Setting.SpeakerID, float64(item.Setting.SpeakerSpeed)/100.0)
		if err != nil {
			// エラーまたはスキップ(キャンセル)された場合は次へ
			continue
		}

		// 2. 接続の生存確認
		vc := p.GetVC()
		if vc == nil || vc.ChannelID() == nil {
			continue
		}

		// 3. Providerをセットする直前に、Speak状態にして少し待つ
		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagMicrophone)

		p.vcMu.Lock()
		p.initFrames = 0 // カウンタをリセット
		p.playing = true
		p.vcMu.Unlock()

		// ここで登録
		vc.SetOpusFrameProvider(p)

		// 4. 再生
		p.streamAudio(cCtx, audio) // 保護されたローカル変数 cCtx を使う

		// 5. 終了
		p.vcMu.Lock()
		p.playing = false
		p.currentCancel = nil // 使い終わったら安全のためにクリア
		p.vcMu.Unlock()

		vc.SetOpusFrameProvider(nil)
		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)
		cCancel() // リソースリークを防ぐため確実にキャンセルを呼ぶ

		select {
		case <-p.Stop:
			return
		default:
		}
	}
}

func (p *VoicePlayer) streamAudio(ctx context.Context, wav []byte) {
	// 1. ffmpeg
	cmd := exec.Command("ffmpeg",
		"-loglevel", "quiet",
		"-i", "pipe:0", // ファイルではなく標準入力から
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)

	stdin, _ := cmd.StdinPipe()   // 書き込み用パイプ
	stdout, _ := cmd.StdoutPipe() // 読み取り用パイプ

	if err := cmd.Start(); err != nil {
		return
	}

	// メモリ上の wav を FFmpeg に書き込むゴルーチン
	go func() {
		defer stdin.Close()
		stdin.Write(wav)
	}()

	defer cmd.Process.Kill()

	// 2. エンコーダの事前生成 (NewVoicePlayerで作って再利用するのが理想)
	pcm := make([]int16, 960*2)
	byteBuf := make([]byte, len(pcm)*2)

	for {
		// 1. まずスキップや終了をチェック
		select {
		case <-ctx.Done():
			return
		case <-p.Skip:
			return
		default:
			// チェックを抜けたら次へ
		}

		_, err := io.ReadFull(stdout, byteBuf)
		if err != nil {
			break
		}

		for i := 0; i < len(pcm); i++ {
			pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
		}

		opusBuf := make([]byte, 4000)
		n, _ := p.encoder.Encode(pcm, opusBuf)

		// 2. 送信。ここでも Skip を同時に待つ
		select {
		case <-ctx.Done():
			return
		case <-p.Skip:
			return
		case p.currentOpusChan <- opusBuf[:n]:
			// 送信完了
		}
	}
	// 押し出し用の無音フレーム
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
	p.vcMu.RLock()
	cancel := p.currentCancel
	p.vcMu.RUnlock()

	if cancel != nil {
		cancel()
	}

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
