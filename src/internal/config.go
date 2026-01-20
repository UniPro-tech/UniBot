package internal

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

type Colors struct {
	Primary int
	Success int
	Warning int
	Error   int
}

type Config struct {
	BotName       string
	Description   string
	BotVersion    string
	Contributors  []Contributors
	URL           string
	GitHub        string
	Colors        Colors
	SupportServer string
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

var (
	Version   = "latest"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

// envが設定されていない場合のデフォルト値
var (
	BotName       = "UniBot"
	Description   = "UniBotはデジタル創作サークルUniProjectの内製Discord Botです。"
	GitHubRepo    = "UniPro-tech/UniBot"
	HomePage      = "https://unibot.uniproject.jp"
	SupportServer = "https://discord.gg/HYWB2aztr8"
)

func LoadConfig() *Config {
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

	// GitHub APIからコントリビューターを取得する
	res, err := http.Get("https://api.github.com/repos/" + GitHubRepoEnv + "/contributors")
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	var bodyJson []GitHubContributorsResponse
	err = json.Unmarshal(body, &bodyJson)
	if err != nil {
		log.Fatal(err)
	}

	contributors := []Contributors{}
	for _, contributor := range bodyJson {
		contributors = append(contributors, Contributors{
			Username: contributor.Login,
			Profile:  contributor.HTMLURL,
			IsBot:    contributor.UserType == "Bot",
		})
	}

	return &Config{
		BotName:      BotNameEnv,
		Description:  DescriptionEnv,
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
		SupportServer: SupportServerEnv,
	}
}
