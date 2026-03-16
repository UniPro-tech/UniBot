package messageComponent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

const (
	rolePanelOptionValueSeparator = "|"
	rolePanelOptionKeyBytes       = 8
	rolePanelPendingAddTTL        = 10 * time.Minute
)

var errInvalidRolePanelOptionValue = errors.New("invalid rolepanel option value")

type RolePanelPendingAdd struct {
	UserID      string
	GuildID     string
	RoleID      string
	Label       string
	Description string
	Emoji       string
	ExpiresAt   time.Time
}

var rolePanelPendingAdds = struct {
	mu    sync.Mutex
	items map[string]RolePanelPendingAdd
}{
	items: make(map[string]RolePanelPendingAdd),
}

func init() {
	RegisterHandler("rolepanel_select_", HandleRolePanelSelect)
	RegisterHandler("rolepanel_add_", HandleRolePanelAdd)
	RegisterHandler("rolepanel_delete", HandleRolePanelDelete)
	RegisterHandler("rolepanel_remove_", HandleRolePanelRemove)
}

// HandleRolePanelSelect はロールパネルのセレクトメニューを処理します
func HandleRolePanelSelect(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if i.GuildID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドはサーバー内でのみ使用できます。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	customID := i.MessageComponentData().CustomID
	selectedValues := i.MessageComponentData().Values
	messageID := strings.TrimPrefix(customID, "rolepanel_select_")

	repo := repository.NewRolePanelRepository(ctx.DB)
	panel, err := repo.GetByMessageID(messageID)
	if err != nil || panel == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルが見つかりませんでした。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if panel.GuildID != i.GuildID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このパネルは別のサーバーのものです。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	messageRoleIDs, err := RolePanelRoleIDsByOptionKey(i.Message)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルのロール情報の読み取りに失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	panelRoleIDs := panelRoleIDsByOptionKey(panel, messageRoleIDs)
	selectedRoleMap := make(map[string]bool)
	for _, value := range selectedValues {
		_, roleID, err := DecodeRolePanelOptionValue(value)
		if err != nil {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "選択されたロール情報の読み取りに失敗しました。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
		selectedRoleMap[roleID] = true
	}

	member, err := s.GuildMember(i.GuildID, i.Member.User.ID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "メンバー情報の取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Flags: discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	currentRoles := make(map[string]bool)
	for _, roleID := range member.Roles {
		currentRoles[roleID] = true
	}

	var addedRoles []string
	var removedRoles []string
	for _, roleID := range panelRoleIDs {
		hasRole := currentRoles[roleID]
		shouldHaveRole := selectedRoleMap[roleID]

		if shouldHaveRole && !hasRole {
			if err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID); err != nil {
				fmt.Printf("failed to add role %s to user %s in guild %s: %v\n", roleID, i.Member.User.ID, i.GuildID, err)
				continue
			}
			addedRoles = append(addedRoles, roleID)
			continue
		}

		if !shouldHaveRole && hasRole {
			if err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID); err != nil {
				fmt.Printf("failed to remove role %s from user %s in guild %s: %v\n", roleID, i.Member.User.ID, i.GuildID, err)
				continue
			}
			removedRoles = append(removedRoles, roleID)
		}
	}

	newRoles := make(map[string]bool)
	for roleID := range currentRoles {
		newRoles[roleID] = true
	}
	for _, roleID := range addedRoles {
		newRoles[roleID] = true
	}
	for _, roleID := range removedRoles {
		delete(newRoles, roleID)
	}

	var description string
	if len(addedRoles) == 0 && len(removedRoles) == 0 {
		description = "ロールに変更はありませんでした。"
	} else {
		if len(addedRoles) > 0 {
			description += "**追加されたロール:**\n"
			for _, roleID := range addedRoles {
				description += fmt.Sprintf("- <@&%s>\n", roleID)
			}
		}
		if len(removedRoles) > 0 {
			if len(addedRoles) > 0 {
				description += "\n"
			}
			description += "**削除されたロール:**\n"
			for _, roleID := range removedRoles {
				description += fmt.Sprintf("- <@&%s>\n", roleID)
			}
		}
	}

	description += "\n**現在のロール:**\n"
	hasAnyRole := false
	for _, opt := range panel.Options {
		roleID, ok := messageRoleIDs[opt.OptionKey]
		if !ok || !newRoles[roleID] {
			continue
		}
		description += fmt.Sprintf("- <@&%s>\n", roleID)
		hasAnyRole = true
	}
	if !hasAnyRole {
		description += "なし\n"
	}

	selectOptions := buildPanelSelectOptions(panel, messageRoleIDs, newRoles)
	components := buildPanelComponents(panel.MessageID, selectOptions)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールを更新しました",
					Description: description,
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}

func intPtr(i int) *int {
	return &i
}

// HandleRolePanelAdd はロール追加用のパネル選択を処理します
func HandleRolePanelAdd(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if i.GuildID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドはサーバー内でのみ使用できます。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if i.Member.Permissions&discordgo.PermissionManageRoles == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "この操作には「ロールの管理」権限が必要です。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	customID := i.MessageComponentData().CustomID
	values := i.MessageComponentData().Values
	if len(values) == 0 {
		return
	}

	messageID := values[0]
	token := strings.TrimPrefix(customID, "rolepanel_add_")
	pendingAdd, ok := ConsumeRolePanelPendingAdd(token, i.Member.User.ID, i.GuildID)
	if !ok {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ロール追加の有効期限が切れたか、無効な操作です。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	repo := repository.NewRolePanelRepository(ctx.DB)
	panel, err := repo.GetByMessageID(messageID)
	if err != nil || panel == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルが見つかりませんでした。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if panel.GuildID != i.GuildID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このパネルは別のサーバーのものです。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if len(panel.Options) >= 25 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このパネルには最大25個までのロールしか追加できません。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	message, err := s.ChannelMessage(panel.ChannelID, panel.MessageID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルメッセージの取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	messageRoleIDs, err := RolePanelRoleIDsByOptionKey(message)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルのロール情報の読み取りに失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	for _, roleID := range messageRoleIDs {
		if roleID == pendingAdd.RoleID {
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Embeds: []*discordgo.MessageEmbed{
						{
							Title:       "エラー",
							Description: "このロールはすでにパネルに追加されています。",
							Color:       config.Colors.Error,
							Timestamp:   time.Now().Format(time.RFC3339),
						},
					},
					Components: []discordgo.MessageComponent{},
				},
			})
			return
		}
	}

	optionKey, err := NewRolePanelOptionKey()
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ロール情報の準備中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	option := &model.RolePanelOption{
		RolePanelID: panel.ID,
		OptionKey:   optionKey,
		Label:       pendingAdd.Label,
		Description: pendingAdd.Description,
		Emoji:       pendingAdd.Emoji,
	}

	if err := repo.AddOption(option); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ロールの追加中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	panel, _ = repo.GetByMessageID(messageID)
	if err := UpdatePanelMessage(s, panel, config, map[string]string{
		optionKey: pendingAdd.RoleID,
	}); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルメッセージの更新に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールを追加しました",
					Description: fmt.Sprintf("ロール <@&%s> を **%s** に追加しました。", pendingAdd.RoleID, panel.Title),
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}

// HandleRolePanelDelete はパネル削除用のセレクトを処理します
func HandleRolePanelDelete(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if i.GuildID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドはサーバー内でのみ使用できます。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if i.Member.Permissions&discordgo.PermissionManageRoles == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "この操作には「ロールの管理」権限が必要です。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	values := i.MessageComponentData().Values
	if len(values) == 0 {
		return
	}

	messageID := values[0]
	repo := repository.NewRolePanelRepository(ctx.DB)
	panel, err := repo.GetByMessageID(messageID)
	if err != nil || panel == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルが見つかりませんでした。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if panel.GuildID != i.GuildID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このパネルは別のサーバーのものです。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	panelTitle := panel.Title
	_ = s.ChannelMessageDelete(panel.ChannelID, panel.MessageID)

	if err := repo.DeleteByID(panel.ID); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルの削除中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールパネルを削除しました",
					Description: fmt.Sprintf("**%s** を削除しました。", panelTitle),
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}

// HandleRolePanelRemove はロール削除用のパネル選択を処理します
func HandleRolePanelRemove(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	if i.GuildID == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このコマンドはサーバー内でのみ使用できます。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if i.Member.Permissions&discordgo.PermissionManageRoles == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "この操作には「ロールの管理」権限が必要です。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	customID := i.MessageComponentData().CustomID
	values := i.MessageComponentData().Values
	if len(values) == 0 {
		return
	}

	messageID := values[0]
	roleID := strings.TrimPrefix(customID, "rolepanel_remove_")
	repo := repository.NewRolePanelRepository(ctx.DB)
	panel, err := repo.GetByMessageID(messageID)
	if err != nil || panel == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルが見つかりませんでした。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if panel.GuildID != i.GuildID {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このパネルは別のサーバーのものです。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	message, err := s.ChannelMessage(panel.ChannelID, panel.MessageID)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルメッセージの取得に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	messageRoleIDs, err := RolePanelRoleIDsByOptionKey(message)
	if err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルのロール情報の読み取りに失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	var optionKey string
	for key, currentRoleID := range messageRoleIDs {
		if currentRoleID == roleID {
			optionKey = key
			break
		}
	}

	if optionKey == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このロールはパネルに追加されていません。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	option := panelOptionByKey(panel, optionKey)
	if option == nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "このロールの設定が見つかりませんでした。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	if err := repo.DeleteOptionByID(option.ID); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "ロールの削除中にエラーが発生しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	panel, _ = repo.GetByMessageID(messageID)
	if err := UpdatePanelMessage(s, panel, config, nil); err != nil {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "エラー",
						Description: "パネルメッセージの更新に失敗しました。",
						Color:       config.Colors.Error,
						Timestamp:   time.Now().Format(time.RFC3339),
					},
				},
				Components: []discordgo.MessageComponent{},
			},
		})
		return
	}

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールを削除しました",
					Description: fmt.Sprintf("ロール <@&%s> を **%s** から削除しました。", roleID, panel.Title),
					Color:       config.Colors.Success,
					Timestamp:   time.Now().Format(time.RFC3339),
				},
			},
			Components: []discordgo.MessageComponent{},
		},
	})
}

// UpdatePanelMessage はパネルのメッセージを更新します
func UpdatePanelMessage(s *discordgo.Session, panel *model.RolePanel, config *internal.Config, extraRoleIDs map[string]string) error {
	if panel == nil {
		return fmt.Errorf("role panel is nil")
	}

	roleIDsByKey := make(map[string]string)
	message, err := s.ChannelMessage(panel.ChannelID, panel.MessageID)
	if err != nil {
		if len(panel.Options) > 0 && len(extraRoleIDs) == 0 {
			return err
		}
	} else {
		roleIDsByKey, err = RolePanelRoleIDsByOptionKey(message)
		if err != nil {
			return err
		}
	}

	for optionKey, roleID := range extraRoleIDs {
		roleIDsByKey[optionKey] = roleID
	}

	embed := &discordgo.MessageEmbed{
		Title:       panel.Title,
		Description: panel.Description,
		Color:       config.Colors.Primary,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ロールを選択してください",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	selectOptions := buildPanelSelectOptions(panel, roleIDsByKey, nil)
	components := buildPanelComponents(panel.MessageID, selectOptions)

	_, err = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    panel.ChannelID,
		ID:         panel.MessageID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
	return err
}

func buildPanelComponents(messageID string, selectOptions []discordgo.SelectMenuOption) []discordgo.MessageComponent {
	components := []discordgo.MessageComponent{}
	if len(selectOptions) == 0 {
		return components
	}

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    "rolepanel_select_" + messageID,
					Placeholder: "ロールを選択...",
					MinValues:   intPtr(0),
					MaxValues:   len(selectOptions),
					Options:     selectOptions,
				},
			},
		},
	}
}

func buildPanelSelectOptions(panel *model.RolePanel, roleIDsByKey map[string]string, defaultRoles map[string]bool) []discordgo.SelectMenuOption {
	var selectOptions []discordgo.SelectMenuOption
	for _, opt := range panel.Options {
		roleID, ok := roleIDsByKey[opt.OptionKey]
		if !ok {
			continue
		}

		option := discordgo.SelectMenuOption{
			Label:       opt.Label,
			Value:       EncodeRolePanelOptionValue(opt.OptionKey, roleID),
			Description: opt.Description,
			Default:     defaultRoles != nil && defaultRoles[roleID],
		}
		if opt.Emoji != "" {
			option.Emoji = &discordgo.ComponentEmoji{Name: opt.Emoji}
		}
		selectOptions = append(selectOptions, option)
	}
	return selectOptions
}

func panelRoleIDsByOptionKey(panel *model.RolePanel, roleIDsByKey map[string]string) []string {
	var roleIDs []string
	for _, opt := range panel.Options {
		roleID, ok := roleIDsByKey[opt.OptionKey]
		if ok {
			roleIDs = append(roleIDs, roleID)
		}
	}
	return roleIDs
}

func panelOptionByKey(panel *model.RolePanel, optionKey string) *model.RolePanelOption {
	for _, opt := range panel.Options {
		if opt.OptionKey == optionKey {
			return opt
		}
	}
	return nil
}

func NewRolePanelOptionKey() (string, error) {
	buf := make([]byte, rolePanelOptionKeyBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}

func EncodeRolePanelOptionValue(optionKey, roleID string) string {
	return optionKey + rolePanelOptionValueSeparator + roleID
}

func DecodeRolePanelOptionValue(value string) (string, string, error) {
	parts := strings.SplitN(value, rolePanelOptionValueSeparator, 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("%w: %q", errInvalidRolePanelOptionValue, value)
	}
	return parts[0], parts[1], nil
}

func RolePanelRoleIDsByOptionKey(message *discordgo.Message) (map[string]string, error) {
	roleIDs := make(map[string]string)
	if message == nil {
		return roleIDs, nil
	}

	for _, component := range message.Components {
		row, ok := component.(discordgo.ActionsRow)
		if !ok {
			continue
		}

		for _, child := range row.Components {
			menu, ok := child.(discordgo.SelectMenu)
			if !ok {
				continue
			}

			for _, option := range menu.Options {
				optionKey, roleID, err := DecodeRolePanelOptionValue(option.Value)
				if err != nil {
					return nil, err
				}
				roleIDs[optionKey] = roleID
			}
		}
	}

	return roleIDs, nil
}

func SaveRolePanelPendingAdd(add RolePanelPendingAdd) (string, error) {
	token, err := newRolePanelPendingAddToken()
	if err != nil {
		return "", err
	}

	add.ExpiresAt = time.Now().Add(rolePanelPendingAddTTL)

	rolePanelPendingAdds.mu.Lock()
	defer rolePanelPendingAdds.mu.Unlock()

	deleteExpiredRolePanelPendingAddsLocked(time.Now())
	rolePanelPendingAdds.items[token] = add

	return token, nil
}

func ConsumeRolePanelPendingAdd(token, userID, guildID string) (RolePanelPendingAdd, bool) {
	now := time.Now()

	rolePanelPendingAdds.mu.Lock()
	defer rolePanelPendingAdds.mu.Unlock()

	deleteExpiredRolePanelPendingAddsLocked(now)

	add, ok := rolePanelPendingAdds.items[token]
	if !ok {
		return RolePanelPendingAdd{}, false
	}
	if add.UserID != userID || add.GuildID != guildID || now.After(add.ExpiresAt) {
		delete(rolePanelPendingAdds.items, token)
		return RolePanelPendingAdd{}, false
	}

	delete(rolePanelPendingAdds.items, token)
	return add, true
}

func deleteExpiredRolePanelPendingAddsLocked(now time.Time) {
	for token, add := range rolePanelPendingAdds.items {
		if now.After(add.ExpiresAt) {
			delete(rolePanelPendingAdds.items, token)
		}
	}
}

func newRolePanelPendingAddToken() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
