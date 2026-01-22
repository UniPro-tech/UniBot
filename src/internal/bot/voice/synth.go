package voice

import (
	"context"
	"log"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func SynthesizeAndPlay(
	ctx *internal.BotContext,
	s *discordgo.Session,
	personalSetting model.TTSPersonalSetting,
	guildID string,
	text string,
) {
	// ---- synthesize ----
	audio, err := ctx.VoiceVox.Synthesize(
		context.Background(),
		text,
		personalSetting.SpeakerID,
		float64(personalSetting.SpeakerPitch),
	)
	if err != nil {
		log.Println(err)
		return
	}

	vc := s.VoiceConnections[guildID]
	if vc == nil {
		log.Printf("No voice connection for guild %s", guildID)

		repository.NewTTSConnectionRepository(ctx.DB).
			DeleteByGuildID(guildID)

		return
	}

	if !vc.Ready {
		log.Println("VC not ready")
		return
	}

	// ---- play ----
	err = PlayWavBytes(vc, audio)
	if err != nil {
		log.Println("play error:", err)
	}
}
