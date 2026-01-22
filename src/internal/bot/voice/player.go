package voice

import (
	"encoding/binary"
	"io"
	"os"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"github.com/hraban/opus"
)

const (
	frameSize = 960 // 20ms @ 48kHz
	channels  = 2
)

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

	for {
		_, err := io.ReadFull(stdout, byteBuf)
		if err != nil {
			break
		}

		// bytes → int16
		for i := 0; i < len(pcm); i++ {
			pcm[i] = int16(binary.LittleEndian.Uint16(byteBuf[i*2:]))
		}

		opusBuf := make([]byte, 4000)

		n, err := enc.Encode(pcm, opusBuf)
		if err != nil {
			return err
		}

		vc.OpusSend <- opusBuf[:n]
	}

	return cmd.Wait()
}
