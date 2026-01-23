package voice

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"github.com/hraban/opus"
)

type VoicePlayer struct {
	GuildID string
	VC      *discordgo.VoiceConnection

	Queue chan [][]byte
	Stop  chan struct{}
	Skip  chan struct{}
}

func NewVoicePlayer(guildID string, vc *discordgo.VoiceConnection) *VoicePlayer {
	p := &VoicePlayer{
		GuildID: guildID,
		VC:      vc,
		Queue:   make(chan [][]byte, 100),
		Stop:    make(chan struct{}),
		Skip:    make(chan struct{}),
	}

	go p.worker()
	return p
}

func (p *VoicePlayer) worker() {
	for {
		select {
		case frames := <-p.Queue:

		FRAME_LOOP:
			for _, frame := range frames {
				select {
				case <-p.Skip:
					log.Println("skip current audio:", p.GuildID)
					break FRAME_LOOP

				case p.VC.OpusSend <- frame:
				}
			}

		case <-p.Stop:
			log.Println("voice worker stopped:", p.GuildID)
			return
		}
	}
}

func (p *VoicePlayer) SkipCurrent() {
	select {
	case p.Skip <- struct{}{}:
	default:
	}
}

func (p *VoicePlayer) Enqueue(frames [][]byte) {
	select {
	case p.Queue <- frames:
	default:
		log.Println("voice queue full:", p.GuildID)
	}
}

func (p *VoicePlayer) Close() {
	close(p.Stop)
}

const (
	frameSize = 960 // 20ms @ 48kHz
	channels  = 2
)

// wavデータをDiscordの音声チャネルで再生する
func PlayWavBytes(vc *discordgo.VoiceConnection, wav []byte) error {

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

	enc, err := opus.NewEncoder(
		48000,
		channels,
		opus.Application(opus.AppAudio),
	)

	if err != nil {
		return err
	}

	pcm := make([]int16, frameSize*channels)
	byteBuf := make([]byte, len(pcm)*2)

	var frames [][]byte

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

		frame := make([]byte, n)
		copy(frame, opusBuf[:n])
		frames = append(frames, frame)
	}

	player := GetManager().Get(vc.GuildID)
	if player == nil {
		return nil
	}

	player.Enqueue(frames)

	return cmd.Wait()
}
