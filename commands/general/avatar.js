const { SlashCommandBuilder, ChannelType } = require("discord.js");
const Discord = require("discord.js");

module.exports = {
  guildOnly: true, // サーバー専用コマンドかどうか
  data: new SlashCommandBuilder() // スラッシュコマンド登録のため
    .setName("avatar")
    .setDescription("Show Avatar(Beta)")
    /*.addUserOption((option) =>
      option
        // optionの名前
        .setName("Users(Not Enable yet)")
        // optionの説明
        .setDescription("User")
        // optionが必須かどうか
        .setRequired(false)
    )*/,

  async execute(i, client) {
    await i.reply(i.author.avatarURL());
    return `No Data`;
  },
};