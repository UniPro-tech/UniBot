// 必要なやつインポート
const { REST, Routes } = require("discord.js");

// ここに自分の情報いれてね
const token = process.env.DISCORD_TOKEN; // Botのトークン
const clientId = process.env.DISCORD_CLIENT_ID; // BotのID
const guildId = process.env.DISCORD_GUILD_ID; // サーバー内コマンドを確認したいならこれも使う

const rest = new REST({ version: "10" }).setToken(token);

// ギルドコマンドかグローバルコマンドか、片方使ってね！
const fetchCommands = async () => {
  try {
    const commands = await rest.get(
      Routes.applicationCommands(clientId)
      // Routes.applicationCommands(clientId) ← グローバルにしたいならこっちにして
    );

    console.log(`登録されてるコマンド（+サブコマンド/オプション）:`);

    for (const command of commands) {
      console.log(`\n/${command.name}（ID: ${command.id}）: ${command.description}`);

      if (command.options?.length > 0) {
        for (const opt of command.options) {
          if (opt.type === 1) {
            // サブコマンド
            console.log(`  ┗ サブコマンド: ${opt.name} - ${opt.description}`);
            if (opt.options?.length > 0) {
              for (const subOpt of opt.options) {
                console.log(`      ┗ 引数: ${subOpt.name} - ${subOpt.description}`);
              }
            }
          } else if (opt.type === 2) {
            // サブコマンドグループ
            console.log(`  ┣ サブコマンドグループ: ${opt.name} - ${opt.description}`);
            for (const sub of opt.options || []) {
              console.log(`      ┗ サブコマンド: ${sub.name} - ${sub.description}`);
              for (const arg of sub.options || []) {
                console.log(`          ┗ 引数: ${arg.name} - ${arg.description}`);
              }
            }
          } else {
            // 普通のオプション
            console.log(`  ┗ 引数: ${opt.name} - ${opt.description}`);
          }
        }
      }
    }
  } catch (err) {
    console.error("エラー:", err);
  }
};

fetchCommands();
