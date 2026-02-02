package general

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"strconv"
	"strings"
	"time"
	"unibot/internal"

	"github.com/bwmarrin/discordgo"
)

const colorCodeImageName = "color.png"

func LoadColorCodeCommandContext() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "colorcode",
		Description: "カラーコードの画像を表示します",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "code",
				Description: "#RRGGBB形式のカラーコードを指定します",
				Required:    true,
			},
		},
	}
}

func ColorCode(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config
	done := make(chan struct{})
	go func() {
		select {
		case <-done:
			return
		case <-time.After(3 * time.Minute):
			_, err := s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
				Embeds: &[]*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "カラーコードの表示に失敗しました。",
						Color:       config.Colors.Error,
						Footer: &discordgo.MessageEmbedFooter{
							Text:    "Requested by " + i.Member.DisplayName(),
							IconURL: i.Member.AvatarURL(""),
						},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			})
			if err != nil {
				log.Println("Failed to edit deferred interaction on timeout:", err)
			}
		}
	}()
	defer close(done)

	codeOption := ""
	options := i.ApplicationCommandData().Options
	if len(options) > 0 {
		codeOption = options[0].StringValue()
	}

	red, green, blue, hex, err := parseHexColor(codeOption)
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "無効なカラーコードです。例: `#FFAA00` または `FFAA00`",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	img := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: red, G: green, B: blue, A: 255}}, image.Point{}, draw.Src)

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "画像の生成に失敗しました。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + i.Member.DisplayName(),
						IconURL: i.Member.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}

	user := i.User
	if i.Member != nil {
		user = i.Member.User
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Color Code",
		Description: fmt.Sprintf("`%s`", hex),
		Color:       int(red)<<16 | int(green)<<8 | int(blue),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "RGB",
				Value: fmt.Sprintf("%d, %d, %d", red, green, blue),
			},
			{
				Name:  "CMYK",
				Value: formatCMYK(red, green, blue),
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: "attachment://" + colorCodeImageName,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + user.Username,
			IconURL: user.AvatarURL(""),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed},
		Files: []*discordgo.File{
			{
				Name:   colorCodeImageName,
				Reader: bytes.NewReader(buf.Bytes()),
			},
		},
	})
	if err != nil {
		s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{
				{
					Title:       "エラー",
					Description: "カラーコードを表示できませんでした。",
					Color:       config.Colors.Error,
					Footer: &discordgo.MessageEmbedFooter{
						Text:    "Requested by " + user.Username,
						IconURL: user.AvatarURL(""),
					},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
			Flags: discordgo.MessageFlagsEphemeral,
		})
		return
	}
}

func parseHexColor(input string) (uint8, uint8, uint8, string, error) {
	trimmed := strings.TrimSpace(input)
	trimmed = strings.TrimPrefix(trimmed, "#")
	if len(trimmed) == 3 {
		trimmed = fmt.Sprintf("%c%c%c%c%c%c", trimmed[0], trimmed[0], trimmed[1], trimmed[1], trimmed[2], trimmed[2])
	}
	if len(trimmed) != 6 {
		return 0, 0, 0, "", fmt.Errorf("invalid length")
	}
	value, err := strconv.ParseUint(trimmed, 16, 32)
	if err != nil {
		return 0, 0, 0, "", err
	}
	red := uint8(value >> 16)
	green := uint8((value >> 8) & 0xFF)
	blue := uint8(value & 0xFF)
	hex := fmt.Sprintf("#%02X%02X%02X", red, green, blue)
	return red, green, blue, hex, nil
}

func formatCMYK(red, green, blue uint8) string {
	r := float64(red) / 255.0
	g := float64(green) / 255.0
	b := float64(blue) / 255.0

	k := 1.0 - maxFloat64(r, g, b)
	if k >= 1.0 {
		return "0, 0, 0, 100"
	}

	c := (1.0 - r - k) / (1.0 - k)
	m := (1.0 - g - k) / (1.0 - k)
	y := (1.0 - b - k) / (1.0 - k)

	return fmt.Sprintf("%.0f, %.0f, %.0f, %.0f", c*100, m*100, y*100, k*100)
}

func maxFloat64(values ...float64) float64 {
	max := values[0]
	for _, value := range values[1:] {
		if value > max {
			max = value
		}
	}
	return max
}
