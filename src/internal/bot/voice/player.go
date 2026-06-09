package voice

import (
	"context"
	"encoding/binary"
	"io"
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
	GuildID   string
	ChannelID string
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

func NewVoicePlayer(guildID string, channelID string, vc voice.Conn, ctx *internal.BotContext) *VoicePlayer {
	enc, _ := opus.NewEncoder(48000, 2, opus.AppAudio)

	p := &VoicePlayer{
		GuildID:   guildID,
		ChannelID: channelID,
		VC:        vc,
		TextQueue: make(chan QueueItem, 50),
		Stop:      make(chan struct{}),
		opusChan:  make(chan []byte, 10), // バッファを持たせてffmpegのブロックを防ぐ
		encoder:   enc,
	}

	// プロバイダーは一度セットしたら外さない (nilを返せばDisgo側が無音として扱う)
	if vc != nil {
		vc.SetOpusFrameProvider(p)
	}

	go p.worker(ctx)
	return p
}

// 20ms毎にDisgoから呼ばれる超高頻度関数。Mutexロックを排除して高速化。
func (p *VoicePlayer) ProvideOpusFrame() ([]byte, error) {
	select {
	case frame := <-p.opusChan:
		return frame, nil
	default:
		return nil, nil // フレームがない時はnilを返すことでDisgoが送信を待機する
	}
}

func (p *VoicePlayer) Close() {
	p.closeOnce.Do(func() {
		close(p.Stop) // Stopのみ閉じることで、EnqueueText側でのpanicを防ぐ
	})
}

func (p *VoicePlayer) CanProvide() bool {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()

	vc := p.VC
	if vc == nil || vc.Gateway().Status() != voice.StatusReady {
		return false
	}
	return true
}

func (p *VoicePlayer) SetVC(vc voice.Conn) {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	if p.VC == vc {
		return
	}

	p.VC = vc
	if vc != nil {
		vc.SetOpusFrameProvider(p)
	}
}

func (p *VoicePlayer) worker(ctx *internal.BotContext) {
	for {
		select {
		case <-p.Stop:
			return
		case item := <-p.TextQueue:
			cCtx, cCancel := context.WithCancel(context.Background())

			p.cancelMu.Lock()
			p.cancelFn = cCancel
			p.cancelMu.Unlock()

			audio, err := ctx.VoiceVox.Synthesize(cCtx, item.Text, item.Setting.SpeakerID, float64(item.Setting.SpeakerSpeed)/100.0)
			if err != nil {
				continue
			}

			vc := p.GetVC()
			if vc == nil || vc.ChannelID() == nil {
				continue
			}

			_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagMicrophone)

			p.streamAudio(cCtx, audio)

			// ffmpegの処理が終わってもチャネルにはバッファが残っているため、
			// 全てDiscordに送信し終わるまで待機する (これがないと語尾がプツッと切断される)
		WaitLoop:
			for len(p.opusChan) > 0 {
				select {
				case <-cCtx.Done(): // スキップされた場合は直ちに抜ける
					break WaitLoop
				default:
					time.Sleep(10 * time.Millisecond)
				}
			}

			_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)
			cCancel()

			p.cancelMu.Lock()
			p.cancelFn = nil
			p.cancelMu.Unlock()

			// スキップ等で残ってしまった古いフレームを完全にクリアし、次の音声への混入を防ぐ
			p.clearOpusChan()
		}
	}
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

func (p *VoicePlayer) streamAudio(ctx context.Context, wav []byte) {
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-loglevel", "quiet",
		"-i", "pipe:0",
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		return
	}

	defer cmd.Wait()

	go func() {
		defer stdin.Close()
		stdin.Write(wav)
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
		}
	}

	// 終了時のポップノイズ対策
	silenceFrame := []byte{0xF8, 0xFF, 0xFE}
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			return
		case p.opusChan <- silenceFrame:
		}
	}
}

func (p *VoicePlayer) EnqueueText(item QueueItem) {
	// クローズ済みの場合はキューに入れない
	select {
	case <-p.Stop:
		return
	default:
	}

	select {
	case p.TextQueue <- item:
	default:
		// キューが上限(50)に達している場合は破棄（またはエラーハンドリング）
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
