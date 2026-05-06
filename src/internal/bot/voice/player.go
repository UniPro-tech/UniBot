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

	// Skip チャネルは不要になったため削除

	currentOpusChan chan []byte
	currentCtx      context.Context
	currentCancel   context.CancelFunc
	encoder         *opus.Encoder
}

func NewVoicePlayer(guildID string, channelID string, vc voice.Conn, ctx *internal.BotContext) *VoicePlayer {
	enc, _ := opus.NewEncoder(48000, 2, opus.AppAudio)

	p := &VoicePlayer{
		GuildID:         guildID,
		ChannelID:       channelID,
		VC:              vc,
		TextQueue:       make(chan QueueItem, 50),
		Stop:            make(chan struct{}),
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

func (p *VoicePlayer) Close() {
	p.vcMu.Lock()
	defer p.vcMu.Unlock()

	select {
	case <-p.Stop:
	default:
		close(p.Stop)
		close(p.TextQueue)
	}
}

func (p *VoicePlayer) CanProvide() bool {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()

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
		cCtx, cCancel := context.WithCancel(context.Background())
		p.vcMu.Lock()
		p.currentCtx = cCtx
		p.currentCancel = cCancel
		p.vcMu.Unlock()

		audio, err := ctx.VoiceVox.Synthesize(cCtx, item.Text, item.Setting.SpeakerID, float64(item.Setting.SpeakerSpeed)/100.0)
		if err != nil {
			continue
		}

		vc := p.GetVC()
		if vc == nil || vc.ChannelID() == nil {
			continue
		}

		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagMicrophone)

		p.vcMu.Lock()
		p.initFrames = 0
		p.playing = true
		p.vcMu.Unlock()

		vc.SetOpusFrameProvider(p)

		p.streamAudio(cCtx, audio)

		// 終了処理
		p.vcMu.Lock()
		p.playing = false
		p.currentCancel = nil
		p.vcMu.Unlock()

		vc.SetOpusFrameProvider(nil)
		_ = vc.SetSpeaking(context.Background(), voice.SpeakingFlagNone)
		cCancel()

		// ★重要: スキップ等で残ってしまった古いフレームを破棄して次の音声に混ざらないようにする
		for len(p.currentOpusChan) > 0 {
			<-p.currentOpusChan
		}

		select {
		case <-p.Stop:
			return
		default:
		}
	}
}

func (p *VoicePlayer) streamAudio(ctx context.Context, wav []byte) {
	// ★重要: CommandContext に変更。Contextがキャンセルされた瞬間、FFmpegも即Killされる
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

	// ★重要: ゾンビプロセス対策のため必ず Wait を呼ぶ (Process.Kill は CommandContext が自動でやってくれるため Wait だけで良い)
	defer cmd.Wait()

	go func() {
		defer stdin.Close()
		stdin.Write(wav)
	}()

	pcm := make([]int16, 960*2)
	byteBuf := make([]byte, len(pcm)*2)

	for {
		// ctx キャンセルにより ffmpeg が Kill されると、stdout が閉じられ ReadFull は即座に err を返す
		_, err := io.ReadFull(stdout, byteBuf)
		if err != nil {
			break
		}

		for i := 0; i < len(pcm); i++ {
			pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
		}

		opusBuf := make([]byte, 4000)
		n, _ := p.encoder.Encode(pcm, opusBuf)

		select {
		case <-ctx.Done():
			return
		case p.currentOpusChan <- opusBuf[:n]:
		}
	}

	silenceFrame := []byte{0xF8, 0xFF, 0xFE}
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
		// これを呼ぶだけで、Synthesizeの停止・streamAudioの停止・ffmpegの強制終了が連鎖的に全て行われます
		cancel()
	}
}

func (p *VoicePlayer) GetVC() voice.Conn {
	p.vcMu.RLock()
	defer p.vcMu.RUnlock()
	return p.VC
}
