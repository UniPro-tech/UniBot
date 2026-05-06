package general

import (
	"bytes"
	"encoding/json"
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

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

const colorCodeImageName = "color.png"

func LoadColorCodeCommandContext() discord.SlashCommandCreate {
	return discord.SlashCommandCreate{
		Name:        "colorcode",
		Description: "カラーコードの画像を表示します",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Name:        "code",
				Description: "#RRGGBB形式のカラーコードを指定します",
				Required:    true,
			},
		},
	}
}

func ColorCode(ctx *internal.BotContext) func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
	return func(data discord.SlashCommandInteractionData, e *handler.CommandEvent) error {
		config := ctx.Config

		var codeOption string
		if opt, ok := data.Options["code"]; ok {
			if opt.Type == discord.ApplicationCommandOptionTypeString {
				if err := json.Unmarshal(opt.Value, &codeOption); err != nil {
					log.Println(err)
				}
			}
		}

		red, green, blue, hex, err := parseHexColor(codeOption)
		if err != nil {
			log.Print(err)
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEphemeral(true).WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "無効なカラーコードです。例: `#FFAA00` または `FFAA00`",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    "Requested by " + e.User().Username,
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}))
			return err
		}

		img := image.NewRGBA(image.Rect(0, 0, 512, 512))
		draw.Draw(img, img.Bounds(), &image.Uniform{C: color.RGBA{R: red, G: green, B: blue, A: 255}}, image.Point{}, draw.Src)

		var buf bytes.Buffer
		if err := png.Encode(&buf, img); err != nil {
			_, err := e.Client().Rest.CreateFollowupMessage(e.ApplicationID(), e.Token(), discord.NewMessageCreate().WithEphemeral(true).WithEmbeds(discord.Embed{
				Title:       "エラー",
				Description: "画像の生成に失敗しました。",
				Color:       config.Colors.Error,
				Footer: &discord.EmbedFooter{
					Text:    "Requested by " + e.User().Username,
					IconURL: e.User().EffectiveAvatarURL(),
				},
				Timestamp: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			}))
			return err
		}

		responseEmbed := discord.Embed{
			Title:       "Color Code",
			Description: fmt.Sprintf("`%s`", hex),
			Color:       int(red)<<16 | int(green)<<8 | int(blue),
			Fields: []discord.EmbedField{
				{
					Name:  "RGB",
					Value: fmt.Sprintf("%d, %d, %d", red, green, blue),
				},
				{
					Name:  "CMYK",
					Value: formatCMYK(red, green, blue),
				},
			},
			Image: &discord.EmbedResource{
				URL: "attachment://" + colorCodeImageName,
			},
			Footer: &discord.EmbedFooter{
				Text:    "Requested by " + e.User().Username,
				IconURL: e.User().EffectiveAvatarURL(),
			},
			Timestamp: func() *time.Time {
				t := time.Now()
				return &t
			}(),
		}

		_, err = e.Client().Rest.CreateFollowupMessage(
			e.ApplicationID(),
			e.Token(),
			discord.NewMessageCreate().WithEmbeds(responseEmbed).WithFiles(discord.NewFile(colorCodeImageName, fmt.Sprintf("Color code %s's Image", hex), bytes.NewReader(buf.Bytes()))))
		return err
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
