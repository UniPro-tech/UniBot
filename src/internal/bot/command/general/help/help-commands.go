package help

type HelpCommand struct {
	Name        string
	Description string
	Usage       string
}

var HelpCommands = []HelpCommand{
	{
		Name:        "/help",
		Description: "コマンドのヘルプを表示します。",
		Usage:       "テキストチャンネルで/help <command>",
	},
	{
		Name:        "/ping",
		Description: "スピードテストを行います。",
		Usage:       "テキストチャンネルで/ping",
	},
	{
		Name:        "/about",
		Description: "このボットの情報を表示します。",
		Usage:       "テキストチャンネルで/about",
	},
	{
		Name:        "/colorcode",
		Description: "カラーコードの画像を表示します。",
		Usage:       "テキストチャンネルで/colorcode code:#RRGGBB",
	},
	{
		Name:        "/tts <subcommand>",
		Description: "読み上げにまつわるコマンドです。\nサブコマンド一覧:\n・`join` : ボイスチャンネルに参加します。\n・`leave`: ボイスチャンネルから退出します。\n・`skip` : 現在再生中の音声をスキップします。\n・`dict` : 読み上げ辞書の管理を行います。\n・`set` : 読み上げの設定（速度など）を変更します。",
		Usage:       "基本はVC接続中に使用を推奨します。",
	},
	{
		Name:        "/tts dict <subcommands>",
		Description: "読み上げ辞書の管理を行います。\nサブコマンド一覧:\n・`add <word> <definition> <case_sensitive>` : 読み上げ辞書に単語を追加します。\n・`remove` : 読み上げ辞書から単語を削除します。\n・`list` : 登録されている単語の一覧を表示します。",
		Usage:       "基本はVC接続中に使用を推奨します。",
	},
	{
		Name:        "/tts skip",
		Description: "現在の読み上げをスキップします。",
		Usage:       "VC接続中、音声が読み上げられているときに`/tts skip`",
	},
	{
		Name:        "/tts set speed <value>",
		Description: "TTSの再生速度を設定します。",
		Usage:       "VC接続中に`/tts set speed 120`のように指定します。",
	},
	{
		Name:        "/tts set voice",
		Description: "TTSの話者を選択します。",
		Usage:       "VC接続中に`/tts set voice`で表示される選択肢から話者を選びます。",
	},
}
