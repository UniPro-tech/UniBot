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
	GuildID   string
	VC        *discordgo.VoiceConnection
	TextQueue chan QueueItem
	Stop      chan struct{}
	Skip      chan struct{}
}

const (
	frameSize = 960
	channels  = 2
)

// VoicePlayer を作る
func NewVoicePlayer(guildID string, vc *discordgo.VoiceConnection, ctx *internal.BotContext) *VoicePlayer {
	p := &VoicePlayer{
		GuildID:   guildID,
		VC:        vc,
		TextQueue: make(chan QueueItem, 50),
		Stop:      make(chan struct{}),
		Skip:      make(chan struct{}),
	}

	go p.worker(ctx)
	return p
}

// Worker: TextQueue から順に再生
func (p *VoicePlayer) worker(ctx *internal.BotContext) {
	for {
		select {
		case item := <-p.TextQueue:
			// 再生用のキャンセルコンテキスト作成
			playCtx, cancel := context.WithCancel(context.Background())

			// skip 信号を受けたらキャンセル
			go func() {
				select {
				case <-p.Skip:
					log.Println("skip signal: stopping current playback")
					cancel()
				case <-p.Stop:
					cancel()
				}
			}()

			// 音声合成と再生
			log.Println("synthesizing:", item.Text)
			audio, err := ctx.VoiceVox.Synthesize(
				playCtx,
				item.Text,
				item.Setting.SpeakerID,
				float64(item.Setting.SpeakerPitch)/ 100.0,
			)
			if err != nil {
				log.Println("synth error:", err)
				continue
			}

			log.Print("start playing:", item.Text)
			err = p.playAudio(playCtx, audio)
			if err != nil && err != context.Canceled {
				log.Println("play error:", err)
			}

		case <-p.Stop:
			return
		}
	}
}

// playAudio を context 対応
func (p *VoicePlayer) playAudio(ctx context.Context, wav []byte) error {
	tmp, _ := os.CreateTemp("", "tts-*.wav")
	defer os.Remove(tmp.Name())
	tmp.Write(wav)
	tmp.Close()

	cmd := exec.Command("ffmpeg", "-loglevel", "quiet", "-i", tmp.Name(),
		"-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		log.Println("[DEBUG] playAudio goroutine started")
		enc, _ := opus.NewEncoder(48000, channels, opus.AppAudio)
		pcm := make([]int16, frameSize*channels)
		byteBuf := make([]byte, len(pcm)*2)

		frameCount := 0
		for {
			_, err := io.ReadFull(stdout, byteBuf)
			if err != nil {
				// 読み取りが終了した（ストリームの終端）は正常終了とみなす
				if err == io.EOF || err == io.ErrUnexpectedEOF {
					log.Printf("[DEBUG] stream ended naturally (frames processed: %d)", frameCount)
					_ = cmd.Wait()
					done <- nil
					return
				}
				log.Printf("[DEBUG] read error: %v (frames processed: %d)", err, frameCount)
				done <- err
				return
			}
			for i := 0; i < len(pcm); i++ {
				pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
			}
			opusBuf := make([]byte, 4000)
			n, _ := enc.Encode(pcm, opusBuf)

			frameCount++
			if frameCount%100 == 0 {
				log.Printf("[DEBUG] processed %d frames, opus size: %d bytes", frameCount, n)
			}

			if p.VC != nil {
				select {
				case p.VC.OpusSend <- opusBuf[:n]:
				case <-ctx.Done():
					log.Printf("[DEBUG] context canceled, killing ffmpeg (frames processed: %d)", frameCount)
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
					done <- context.Canceled
					return
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		return context.Canceled
	case err := <-done:
		return err
	}
}

// キューに追加
func (p *VoicePlayer) EnqueueText(item QueueItem) {
	select {
	case p.TextQueue <- item:
	default:
		log.Println("text queue full:", p.GuildID)
	}
}

// Skip
func (p *VoicePlayer) SkipCurrent() {
	select {
	case p.Skip <- struct{}{}:
	default:
	}
}

// Close
func (p *VoicePlayer) Close() {
	select {
	case <-p.Stop:
		// すでに閉じてたら何もしない
	default:
		close(p.Stop)
	}
}
