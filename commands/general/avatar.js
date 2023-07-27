const { SlashCommandBuilder, ChannelType } = require("discord.js");
const Discord = require("discord.js");

module.exports = {
  guildOnly: true, // サーバー専用コマンドかどうか
  data: new SlashCommandBuilder() // スラッシュコマンド登録のため
    .setName("avatar")
    .setDescription("Show Avatar(Beta)")
    .addChannelOption((option) =>
      option
        // optionの名前
        .setName("channel")
        // optionの説明
        .setDescription("The channel to join")
        // optionが必須かどうか
        .setRequired(false)
        .addUserOption((option) =>
          option
            .setName("User")
            .setDescription("SetUser")
        )
    ),

  async execute(i, client) {
    await i.reply(i.author.avatarURL());
    return `No Data`;
  },
};