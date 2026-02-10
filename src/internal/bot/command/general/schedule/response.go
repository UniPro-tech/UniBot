package schedule

import "github.com/bwmarrin/discordgo"

func RespondEdit(s *discordgo.Session, i *discordgo.InteractionCreate, data *discordgo.InteractionResponseData) error {
	if data == nil {
		data = &discordgo.InteractionResponseData{}
	}

	edit := &discordgo.WebhookEdit{Flags: data.Flags}
	if data.Content != "" {
		content := data.Content
		edit.Content = &content
	}
	if data.Components != nil {
		components := data.Components
		edit.Components = &components
	}
	if data.Embeds != nil {
		embeds := data.Embeds
		edit.Embeds = &embeds
	}
	if data.AllowedMentions != nil {
		edit.AllowedMentions = data.AllowedMentions
	}
	if data.Files != nil {
		edit.Files = data.Files
	}
	if data.Attachments != nil {
		edit.Attachments = data.Attachments
	}

	_, err := s.InteractionResponseEdit(i.Interaction, edit)
	return err
}
