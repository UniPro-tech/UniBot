package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"unibot/internal/api/voicevox"

	"gorm.io/gorm"
)

const (
	// contributorsTimeout は GitHub API の呼び出しに許す時間。
	// 起動を無期限に待たせないために設ける。
	contributorsTimeout = 10 * time.Second
	// maxContributorsBody は読み込むレスポンス本文の上限。
	maxContributorsBody = 1 << 20
)

type Colors struct {
	Primary int
	Success int
	Warning int
	Error   int
}

type Config struct {
	BotName        string
	Description    string
	AdminGuildID   string
	AdminRoleID    string
	BotVersion     string
	Contributors   []Contributors
	URL            string
	GitHub         string
	Colors         Colors
	SupportServer  string
	VoiceVoxURI    string
	VoiceVoxAPIKey string
}

type GitHubContributorsResponse struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	UserType          string `json:"type"`
	UserViewType      string `json:"user_view_type"`
	SiteAdmin         bool   `json:"site_admin"`
	Contributions     int    `json:"contributions"`
}

type Contributors struct {
	Username string `json:"login"`
	Profile  string `json:"html_url"`
	IsBot    bool   `json:"site_admin"`
}

type BotContext struct {
	DB       *gorm.DB
	Config   *Config
	VoiceVox *voicevox.Client
}

var (
	Version   = "latest"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

// envが設定されていない場合のデフォルト値
var (
	BotName        = "UniBot"
	Description    = "UniBotはデジタル創作サークルUniProjectの内製Discord Botです。"
	AdminGuildID   = "1191346186880286770"
	AdminRoleID    = "1390633352360628234"
	GitHubRepo     = "UniPro-tech/UniBot"
	HomePage       = "https://unibot.uniproject.jp"
	SupportServer  = "https://discord.gg/HYWB2aztr8"
	VoiceVoxURI    = "http://localhost:50021"
	VoiceVoxAPIKey = ""
)

var (
	configOnce   sync.Once
	cachedConfig *Config
)

// LoadConfig は設定を読み込む。結果はプロセス内で1度だけ構築され、以降は再利用される。
func LoadConfig() *Config {
	configOnce.Do(func() {
		cachedConfig = loadConfig()
	})
	return cachedConfig
}

func loadConfig() *Config {
	version := Version

	if Version == "latest" {
		version = GitBranch + "@" + GitCommit
	} else {
		version = Version + "+" + GitCommit
	}

	// envから設定を読み込む
	BotNameEnv := os.Getenv("CONFIG_BOT_NAME")
	if BotNameEnv == "" {
		BotNameEnv = BotName
	}
	DescriptionEnv := os.Getenv("CONFIG_DESCRIPTION")
	if DescriptionEnv == "" {
		DescriptionEnv = Description
	}
	AdminGuildIDEnv := os.Getenv("CONFIG_ADMIN_GUILD_ID")
	if AdminGuildIDEnv == "" {
		AdminGuildIDEnv = AdminGuildID
	}
	AdminRoleIDEnv := os.Getenv("CONFIG_ADMIN_ROLE_ID")
	if AdminRoleIDEnv == "" {
		AdminRoleIDEnv = AdminRoleID
	}
	GitHubRepoEnv := os.Getenv("CONFIG_GITHUB_REPO")
	if GitHubRepoEnv == "" {
		GitHubRepoEnv = GitHubRepo
	}
	HomePageEnv := os.Getenv("CONFIG_HOME_PAGE")
	if HomePageEnv == "" {
		HomePageEnv = HomePage
	}
	SupportServerEnv := os.Getenv("CONFIG_SUPPORT_SERVER")
	if SupportServerEnv == "" {
		SupportServerEnv = SupportServer
	}
	VoiceVoxURIEnv := os.Getenv("VOICEVOX_URI")
	if VoiceVoxURIEnv == "" {
		VoiceVoxURIEnv = VoiceVoxURI
	}
	VoiceVoxAPIKeyEnv := os.Getenv("VOICEVOX_API_KEY")
	if VoiceVoxAPIKeyEnv == "" {
		VoiceVoxAPIKeyEnv = VoiceVoxAPIKey
	}

	// コントリビューターの取得に失敗しても致命的ではない（/about の表示が欠けるだけ）。
	contributors, err := fetchContributors(context.Background(), GitHubRepoEnv)
	if err != nil {
		slog.Warn("failed to fetch contributors from github",
			slog.String("repo", GitHubRepoEnv), slog.Any("err", err))
	}

	return &Config{
		BotName:      BotNameEnv,
		Description:  DescriptionEnv,
		AdminGuildID: AdminGuildIDEnv,
		AdminRoleID:  AdminRoleIDEnv,
		BotVersion:   version,
		Contributors: contributors,
		URL:          HomePageEnv,
		GitHub:       "https://github.com/" + GitHubRepoEnv,
		Colors: Colors{
			Primary: 0x3498DB,
			Success: 0x2ECC71,
			Warning: 0xF1C40F,
			Error:   0xE74C3C,
		},
		SupportServer:  SupportServerEnv,
		VoiceVoxURI:    VoiceVoxURIEnv,
		VoiceVoxAPIKey: VoiceVoxAPIKeyEnv,
	}
}

// fetchContributors は GitHub API からコントリビューター一覧を取得する。
// 失敗しても呼び出し側が継続できるよう、エラーを返すだけでプロセスは落とさない。
func fetchContributors(ctx context.Context, repo string) ([]Contributors, error) {
	url := "https://api.github.com/repos/" + repo + "/contributors"

	ctx, cancel := context.WithTimeout(ctx, contributorsTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if id, secret := os.Getenv("GITHUB_OAUTH_ID"), os.Getenv("GITHUB_OAUTH_SECRET"); id != "" && secret != "" {
		req.SetBasicAuth(id, secret)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api returned %s", res.Status)
	}

	body, err := io.ReadAll(io.LimitReader(res.Body, maxContributorsBody))
	if err != nil {
		return nil, err
	}

	var response []GitHubContributorsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		// レスポンス本文にはトークン等が含まれうるため、内容は出力しない。
		return nil, fmt.Errorf("failed to decode contributors response: %w", err)
	}

	contributors := make([]Contributors, 0, len(response))
	for _, contributor := range response {
		contributors = append(contributors, Contributors{
			Username: contributor.Login,
			Profile:  contributor.HTMLURL,
			IsBot:    contributor.UserType == "Bot",
		})
	}
	return contributors, nil
}
