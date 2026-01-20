package internal

type Colors struct {
	Primary int
	Success int
	Warning int
	Error   int
}

type Config struct {
	BotName       string
	BotVersion    string
	Contributors  []string
	URL           string
	GitHub        string
	Colors        Colors
	SupportServer string
}

var (
	Version   = "latest"
	GitCommit = "unknown"
	GitBranch = "unknown"
)

func LoadConfig() *Config {
	version := Version

	if Version == "latest" {
		version = GitBranch + "@" + GitCommit
	} else {
		version = Version + "+" + GitCommit
	}

	return &Config{
		BotName:      "UniBot",
		BotVersion:   version,
		Contributors: []string{"Yuito Akatsuki <yuito@yuito-it.jp>"},
		URL:          "https://unibot.uniproject.jp",
		GitHub:       "https://github.com/UniPro-tech/UniBot",
		Colors: Colors{
			Primary: 0x3498DB,
			Success: 0x2ECC71,
			Warning: 0xF1C40F,
			Error:   0xE74C3C,
		},
		SupportServer: "https://discord.gg/your-invite-code",
	}
}
