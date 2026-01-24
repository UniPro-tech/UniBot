package messageComponent

import (
	"fmt"
	"strings"
	"time"
	"unibot/internal"
	"unibot/internal/model"
	"unibot/internal/repository"

	"github.com/bwmarrin/discordgo"
)

func init() {
	RegisterHandler("rolepanel_select_", HandleRolePanelSelect)
	RegisterHandler("rolepanel_add_", HandleRolePanelAdd)
	RegisterHandler("rolepanel_delete", HandleRolePanelDelete)
	RegisterHandler("rolepanel_remove_", HandleRolePanelRemove)
}

// HandleRolePanelSelect はロールパネルのセレクトメニューを処理します
func HandleRolePanelSelect(ctx *internal.BotContext, s *discordgo.Session, i *discordgo.InteractionCreate) {
	config := ctx.Config

	// Guildチェック
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
	selectedRoleIDs := i.MessageComponentData().Values

	// CustomIDからメッセージIDを取得 (rolepanel_select_messageID)
	messageID := strings.TrimPrefix(customID, "rolepanel_select_")

	// パネルを取得
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

	// GuildID一致チェック
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

	// パネルに登録されている全ロールIDを取得
	panelRoleIDs := make(map[string]bool)
	for _, opt := range panel.Options {
		panelRoleIDs[opt.RoleID] = true
	}

	// 選択されたロールをマップに変換
	selectedMap := make(map[string]bool)
	for _, roleID := range selectedRoleIDs {
		selectedMap[roleID] = true
	}

	var addedRoles []string
	var removedRoles []string

	// ユーザーの現在のロールを確認
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

	// ロールの追加・削除を処理
	for roleID := range panelRoleIDs {
		hasRole := currentRoles[roleID]
		shouldHaveRole := selectedMap[roleID]

		if shouldHaveRole && !hasRole {
			// ロールを追加
			err := s.GuildMemberRoleAdd(i.GuildID, i.Member.User.ID, roleID)
			if err != nil {
				fmt.Printf("failed to add role %s to user %s in guild %s: %v\n", roleID, i.Member.User.ID, i.GuildID, err)
				continue
			}
			addedRoles = append(addedRoles, roleID)
		} else if !shouldHaveRole && hasRole {
			// ロールを削除
			err := s.GuildMemberRoleRemove(i.GuildID, i.Member.User.ID, roleID)
			if err != nil {
				fmt.Printf("failed to remove role %s from user %s in guild %s: %v\n", roleID, i.Member.User.ID, i.GuildID, err)
				continue
			}
			removedRoles = append(removedRoles, roleID)
		}
	}

	// 更新後のユーザーのロール状態を計算
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

	// 結果メッセージを作成
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

	// 現在のロール状態を表示
	description += "\n**現在のロール:**\n"
	hasAnyRole := false
	for _, opt := range panel.Options {
		if newRoles[opt.RoleID] {
			description += fmt.Sprintf("- <@&%s>\n", opt.RoleID)
			hasAnyRole = true
		}
	}
	if !hasAnyRole {
		description += "なし\n"
	}

	// ユーザー専用のセレクトメニューを作成（現在のロールがデフォルト選択）
	var selectOptions []discordgo.SelectMenuOption
	for _, opt := range panel.Options {
		option := discordgo.SelectMenuOption{
			Label:       opt.Label,
			Value:       opt.RoleID,
			Description: opt.Description,
			Default:     newRoles[opt.RoleID],
		}
		if opt.Emoji != "" {
			option.Emoji = &discordgo.ComponentEmoji{
				Name: opt.Emoji,
			}
		}
		selectOptions = append(selectOptions, option)
	}

	var components []discordgo.MessageComponent
	if len(selectOptions) > 0 {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "rolepanel_select_" + panel.MessageID,
						Placeholder: "ロールを選択...",
						MinValues:   intPtr(0),
						MaxValues:   len(selectOptions),
						Options:     selectOptions,
					},
				},
			},
		}
	}

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

	// Guildチェック
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

	// 権限チェック
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

	// CustomIDからロール情報をデコード (rolepanel_add_roleID|label|description|emoji)
	parts := strings.TrimPrefix(customID, "rolepanel_add_")
	data := strings.Split(parts, "|")
	if len(data) < 2 {
		return
	}

	roleID := data[0]
	label := data[1]
	description := ""
	emoji := ""
	if len(data) > 2 {
		description = data[2]
	}
	if len(data) > 3 {
		emoji = data[3]
	}

	repo := repository.NewRolePanelRepository(ctx.DB)

	// パネルを取得
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

	// GuildID一致チェック
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

	// オプション数の上限チェック (Discord制限: 25)
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

	// 重複チェック
	for _, opt := range panel.Options {
		if opt.RoleID == roleID {
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

	// オプションを追加
	option := &model.RolePanelOption{
		RolePanelID: panel.ID,
		RoleID:      roleID,
		Label:       label,
		Description: description,
		Emoji:       emoji,
	}

	err = repo.AddOption(option)
	if err != nil {
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

	// パネルを再取得してメッセージを更新
	panel, _ = repo.GetByMessageID(messageID)
	UpdatePanelMessage(s, panel, config)

	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "ロールを追加しました",
					Description: fmt.Sprintf("ロール <@&%s> を **%s** に追加しました。", roleID, panel.Title),
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

	// Guildチェック
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

	// 権限チェック
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

	// パネルを取得
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

	// GuildID一致チェック
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

	// メッセージを削除
	_ = s.ChannelMessageDelete(panel.ChannelID, panel.MessageID)

	// データベースから削除
	err = repo.DeleteByID(panel.ID)
	if err != nil {
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

	// Guildチェック
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

	// 権限チェック
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

	// パネルを取得
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

	// GuildID一致チェック
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

	// ロールを探す
	var optionID uint
	found := false
	for _, opt := range panel.Options {
		if opt.RoleID == roleID {
			optionID = opt.ID
			found = true
			break
		}
	}

	if !found {
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

	// オプションを削除
	err = repo.DeleteOptionByID(optionID)
	if err != nil {
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

	// パネルを再取得してメッセージを更新
	panel, _ = repo.GetByMessageID(messageID)
	UpdatePanelMessage(s, panel, config)

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
func UpdatePanelMessage(s *discordgo.Session, panel *model.RolePanel, config *internal.Config) {
	embed := &discordgo.MessageEmbed{
		Title:       panel.Title,
		Description: panel.Description,
		Color:       config.Colors.Primary,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "ロールを選択してください",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var selectOptions []discordgo.SelectMenuOption
	for _, opt := range panel.Options {
		option := discordgo.SelectMenuOption{
			Label:       opt.Label,
			Value:       opt.RoleID,
			Description: opt.Description,
		}
		if opt.Emoji != "" {
			option.Emoji = &discordgo.ComponentEmoji{
				Name: opt.Emoji,
			}
		}
		selectOptions = append(selectOptions, option)
	}

	var components []discordgo.MessageComponent
	if len(selectOptions) > 0 {
		components = []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "rolepanel_select_" + panel.MessageID,
						Placeholder: "ロールを選択...",
						MinValues:   intPtr(0),
						MaxValues:   len(selectOptions),
						Options:     selectOptions,
					},
				},
			},
		}
	}

	_, _ = s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    panel.ChannelID,
		ID:         panel.MessageID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
}
