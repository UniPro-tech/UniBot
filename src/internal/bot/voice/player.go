package voice

import (
	"context"
	"encoding/binary"
	"io"
	"log/slog"
	"os/exec"
	"sync"
	"time"

	"unibot/internal"
	"unibot/internal/logger"

	"github.com/disgoorg/disgo/voice"
	"github.com/hraban/opus"
)

type QueueItem struct {
	// Ctx は投入元のイベントの trace を引き継ぐためのもの。
	// 読み上げは非同期に処理されるため、http.Request と同様に構造体で持ち回る。
	// nil の場合は context.Background() として扱う。
	Ctx  context.Context
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

func NewVoicePlayer(guildID int64, channelID int64, vc voice.Conn, bctx *internal.BotContext) *VoicePlayer {
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

	go p.worker(bctx)
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

func (p *VoicePlayer) worker(bctx *internal.BotContext) {
	for {
		select {
		case <-p.Stop:
			return
		case item := <-p.TextQueue:
			p.processItem(bctx, item)
		}
	}
}

func (p *VoicePlayer) processItem(bctx *internal.BotContext, item QueueItem) {
	// 投入元の trace を引き継ぎつつ、処理単位として新しい request_id を発行する。
	ctx, _ := logger.WithRequest(queueItemContext(item))

	slog.DebugContext(ctx, "tts item dequeued",
		slog.Int64("guild_id", p.GuildID),
		// 本文は載せず、長さのみ記録する。
		slog.Int("text_len", len([]rune(item.Text))),
	)
	start := time.Now()

	cCtx, cCancel := context.WithCancel(ctx)
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
	audio, err := bctx.VoiceVox.Synthesize(cCtx, item.Text, "0", float64(100)/100.0)
	if err != nil {
		slog.ErrorContext(ctx, "tts synthesis failed",
			slog.Int64("guild_id", p.GuildID), slog.Any("err", err))
		return
	}

	vc := p.GetVC()
	if vc == nil {
		slog.WarnContext(ctx, "tts skipped: no voice connection",
			slog.Int64("guild_id", p.GuildID))
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
			slog.DebugContext(ctx, "tts playback cancelled", slog.Int64("guild_id", p.GuildID))
			break WaitLoop
		case <-maxWait:
			slog.WarnContext(ctx, "tts playback drain timeout",
				slog.Int64("guild_id", p.GuildID), slog.Int("pending_frames", len(p.opusChan)))
			break WaitLoop
		default:
			if !p.CanProvide() {
				slog.WarnContext(ctx, "tts playback aborted: connection lost",
					slog.Int64("guild_id", p.GuildID))
				break WaitLoop
			}
			time.Sleep(10 * time.Millisecond)
		}
	}

	_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)

	slog.DebugContext(ctx, "tts playback finished",
		slog.Int64("guild_id", p.GuildID),
		slog.Int64("duration_ms", time.Since(start).Milliseconds()),
	)
}

// queueItemContext はキュー項目に載っている context を安全に取り出す。
func queueItemContext(item QueueItem) context.Context {
	if item.Ctx == nil {
		return context.Background()
	}
	return item.Ctx
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
		slog.ErrorContext(ctx, "failed to open ffmpeg stdin",
			slog.Int64("guild_id", p.GuildID), slog.Any("err", err))
		return
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.ErrorContext(ctx, "failed to open ffmpeg stdout",
			slog.Int64("guild_id", p.GuildID), slog.Any("err", err))
		return
	}

	if err := cmd.Start(); err != nil {
		slog.ErrorContext(ctx, "failed to start ffmpeg",
			slog.Int64("guild_id", p.GuildID), slog.Any("err", err))
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
	ctx := queueItemContext(item)

	select {
	case <-p.Stop:
		slog.WarnContext(ctx, "tts enqueue on closed player", slog.Int64("guild_id", p.GuildID))
		return
	default:
	}

	select {
	case p.TextQueue <- item:
		slog.DebugContext(ctx, "tts item enqueued",
			slog.Int64("guild_id", p.GuildID), slog.Int("queue_len", len(p.TextQueue)))
	default:
		// キューが上限(50)に達している場合は破棄する。
		slog.WarnContext(ctx, "tts queue full, dropping item",
			slog.Int64("guild_id", p.GuildID), slog.Int("queue_cap", cap(p.TextQueue)))
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
