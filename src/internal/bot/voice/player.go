package voice

import (
	"context"
	"encoding/binary"
	"io"
	"os/exec"
	"sync"
	"time"
	"unibot/internal"

	"github.com/disgoorg/disgo/voice"
	"github.com/hraban/opus"
)

type QueueItem struct {
	Text string
}

type VoicePlayer struct {
	GuildID   int64
	ChannelID int64
	VC        voice.Conn

	TextQueue chan QueueItem
	Stop      chan struct{}

	opusChan chan []byte
	encoder  *opus.Encoder

	vcMu sync.RWMutex

	cancelMu sync.Mutex
	cancelFn context.CancelFunc

	closeOnce sync.Once
}

func NewVoicePlayer(guildID int64, channelID int64, vc voice.Conn, ctx *internal.BotContext) *VoicePlayer {
	enc, _ := opus.NewEncoder(48000, 2, opus.AppAudio)

	p := &VoicePlayer{
		GuildID:   guildID,
		ChannelID: channelID,
		VC:        vc,
		TextQueue: make(chan QueueItem, 50),
		Stop:      make(chan struct{}),
		opusChan:  make(chan []byte, 10),
		encoder:   enc,
	}

	// プロバイダーは一度セットしたら外さない (nilを返せばDisgo側が無音として扱う)
	if vc != nil {
		vc.SetOpusFrameProvider(p)
	}

	go p.worker(ctx)
	return p
}

// 20ms毎にDisgoから呼ばれる関数
func (p *VoicePlayer) ProvideOpusFrame() ([]byte, error) {
	select {
	case frame := <-p.opusChan:
		return frame, nil
	default:
		return nil, nil
	}
}

func (p *VoicePlayer) Close() {
	p.closeOnce.Do(func() {
		close(p.Stop) // Stopのみ閉じることで、EnqueueText側でのpanicを防ぐ
		p.SkipCurrent()
	})
}

// CanProvide の判定を緩和・安全化
func (p *VoicePlayer) CanProvide() bool {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()

	vc := p.VC
	if vc == nil {
		return false
	}

	// ChannelID が存在すれば接続中とみなす
	if vc.ChannelID() == nil {
		return false
	}

	return true
}

func (p *VoicePlayer) SetVC(vc voice.Conn) {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	if vc == nil {
		p.VC = nil
		return
	}

	p.VC = vc
	// 接続インスタンスが変わった場合は必ず Provider を再登録する
	vc.SetOpusFrameProvider(p)
}

func (p *VoicePlayer) worker(ctx *internal.BotContext) {
	for {
		select {
		case <-p.Stop:
			return
		case item := <-p.TextQueue:
			// TODO: log
			p.processItem(ctx, item)
		}
	}
}

func (p *VoicePlayer) processItem(ctx *internal.BotContext, item QueueItem) {
	cCtx, cCancel := context.WithCancel(context.Background())
	defer cCancel()

	p.cancelMu.Lock()
	p.cancelFn = cCancel
	p.cancelMu.Unlock()

	defer func() {
		p.cancelMu.Lock()
		p.cancelFn = nil
		p.cancelMu.Unlock()
		p.clearOpusChan()
	}()

	// 2. 音声合成
	audio, err := ctx.VoiceVox.Synthesize(cCtx, item.Text, "0", float64(100)/100.0)
	if err != nil {
		// TODO: log
		return
	}

	vc := p.GetVC()
	if vc == nil {
		// TODO: log
		return
	}

	_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagMicrophone)
	p.streamAudio(cCtx, audio)

	// 待機ループ
	maxWait := time.After(5 * time.Second)
WaitLoop:
	for len(p.opusChan) > 0 {
		select {
		case <-cCtx.Done():
			// TODO: log
			break WaitLoop
		case <-maxWait:
			// TODO: log
			break WaitLoop
		default:
			if !p.CanProvide() {
				// TODO: log
				break WaitLoop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)
	// TODO: log
}

func (p *VoicePlayer) clearOpusChan() {
	for {
		select {
		case <-p.opusChan:
		default:
			return
		}
	}
}

func (p *VoicePlayer) streamAudio(parentCtx context.Context, wav []byte) {
	// streamAudio 専用のキャンセル可能な Context を作成
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel() // 関数の脱出時に必ず Context をキャンセル（FFmpeg に SIGKILL が送られる）

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-loglevel", "quiet",
		"-i", "pipe:0",
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	// 関数の終了処理を安全な順序で一括実行
	defer func() {
		cancel()                           // 1. FFmpeg プロセスに終了シグナルを送る
		_ = stdin.Close()                  // 2. stdin を閉じる
		_, _ = io.Copy(io.Discard, stdout) // 3. stdout の残りを読み捨ててパイプ詰まりを防止
		_ = cmd.Wait()                     // 4. 安全にプロセスの完全終了を待つ
	}()

	go func() {
		defer stdin.Close()
		_, _ = stdin.Write(wav)
	}()

	pcm := make([]int16, 960*2)
	byteBuf := make([]byte, len(pcm)*2)

	for {
		_, err := io.ReadFull(stdout, byteBuf)
		if err != nil {
			break
		}

		for i := range pcm {
			pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
		}

		opusBuf := make([]byte, 4000)
		n, err := p.encoder.Encode(pcm, opusBuf)
		if err != nil {
			continue
		}

		select {
		case <-ctx.Done():
			return
		case p.opusChan <- opusBuf[:n]:
		case <-time.After(500 * time.Millisecond):
			// タイムアウト時、 return すると defer された cancel() が走り FFmpeg が安全に終了する
			return
		}
	}

	// 終了時のポップノイズ対策
	silenceFrame := []byte{0xF8, 0xFF, 0xFE}
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return
		case p.opusChan <- silenceFrame:
		case <-time.After(100 * time.Millisecond):
			return
		}
	}
}

func (p *VoicePlayer) EnqueueText(item QueueItem) {
	// クローズ済みの場合はキューに入れない
	select {
	case <-p.Stop:
		// TODO: log
		return
	default:
	}

	select {
	case p.TextQueue <- item:
		// TODO: log
	default:
		// キューが上限(50)に達している場合は破棄（またはエラーハンドリング）
		// TODO: log
	}
}

func (p *VoicePlayer) SkipCurrent() {
	p.cancelMu.Lock()
	cancel := p.cancelFn
	p.cancelMu.Unlock()

	if cancel != nil {
		cancel()
	}
}

func (p *VoicePlayer) GetVC() voice.Conn {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()
	return p.VC
}
