package voice

import (
	"context"
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"
	"unibot/internal"
	"unibot/internal/model"

	"github.com/bwmarrin/discordgo"
	"github.com/hraban/opus"
)

type QueueItem struct {
	Text    string
	Setting model.TTSPersonalSetting
}

type VoicePlayer struct {
	GuildID string
	VC      *discordgo.VoiceConnection

	TextQueue chan QueueItem
	Stop      chan struct{}
	Skip      chan struct{}
}

const (
	frameSize = 960 // 20ms @ 48kHz
	channels  = 2
)

func NewVoicePlayer(guildID string, vc *discordgo.VoiceConnection, ctx *internal.BotContext) *VoicePlayer {
	p := &VoicePlayer{
		GuildID:   guildID,
		VC:        vc,
		TextQueue: make(chan QueueItem, 50), // 文章単位で十分
		Stop:      make(chan struct{}),
		Skip:      make(chan struct{}),
	}

	go p.worker(ctx)
	return p
}

// TextQueue から TTS → 再生
func (p *VoicePlayer) worker(ctx *internal.BotContext) {
	for {
		select {
		case item := <-p.TextQueue:
			audio, err := ctx.VoiceVox.Synthesize(
				context.Background(),
				item.Text,
				item.Setting.SpeakerID,
				float64(item.Setting.SpeakerPitch),
			)
			if err != nil {
				log.Println("synth error:", err)
				continue
			}

			err = p.playAudio(audio)
			if err != nil {
				log.Println("play error:", err)
			}

		case <-p.Skip:
			log.Println("skip signal: stopping current playback")
			// 実際に再生中の ffmpeg は kill して次の TextQueue に進む
			// TODO: 管理するなら context.WithCancel() で制御

		case <-p.Stop:
			return
		}
	}
}

// 音声を VC に送信
func (p *VoicePlayer) playAudio(wav []byte) error {
	tmp, err := os.CreateTemp("", "tts-*.wav")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	tmp.Write(wav)
	tmp.Close()

	cmd := exec.Command(
		"ffmpeg",
		"-loglevel", "quiet",
		"-i", tmp.Name(),
		"-f", "s16le",
		"-ar", "48000",
		"-ac", "2",
		"pipe:1",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	enc, err := opus.NewEncoder(48000, channels, opus.AppAudio)
	if err != nil {
		return err
	}

	pcm := make([]int16, frameSize*channels)
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
		n, err := enc.Encode(pcm, opusBuf)
		if err != nil {
			return err
		}

		if p.VC != nil {
			p.VC.OpusSend <- opusBuf[:n]
		}
	}

	return cmd.Wait()
}

// キューに TTS 追加
func (p *VoicePlayer) EnqueueText(item QueueItem) {
	select {
	case p.TextQueue <- item:
	default:
		log.Println("text queue full:", p.GuildID)
	}
}

// skip 再生中の音声をスキップ
func (p *VoicePlayer) SkipCurrent() {
	select {
	case p.Skip <- struct{}{}:
	default:
	}
}

// クローズ
func (p *VoicePlayer) Close() {
	if p.Stop == nil {
		return
	}
	close(p.Stop)
	p.Stop = nil
}
