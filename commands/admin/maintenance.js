const { SlashCommandBuilder, Activity } = require("discord.js");
const Discord = require("discord.js");
module.exports = {
  guildOnly: false, // サーバー専用コマンドかどうか
  adminGuildOnly: true,
  data: new SlashCommandBuilder() // スラッシュコマンド登録のため
    .setName("maintenance")
    .setDescription("メンテモード")
    .addStringOption(option => option.setName('enablet').setDescription('on/off'))
    .addStringOption(option => option.setName('playing').setDescription('プレイ中に表示するやつ'))
    .addStringOption(option => option.setName('status').setDescription('すたーてす')),

  async execute(i, client, command) {
    try {
      const onoff = i.options.getString('enablet');
      if (onoff == 'on') {
        const status = i.options.getString('status');
        const playing = i.options.getString('playing');
        client.user.setActivity(playing);
        client.user.setStatus(status);
        const embed = new Discord.EmbedBuilder()
          .setTitle("ok")
          .setColor(client.config.color.s)
          .setTimestamp();

        i.reply({ embeds: [embed] });
        return `{ "onoff":"on","status": "${status}", "playing": "${playing}" }`;
      } else {
        const embed = new Discord.EmbedBuilder()
          .setTitle("ok")
          .setColor(client.config.color.s)
          .setTimestamp();

        i.reply({ embeds: [embed] });
        return `{ "onoff":"off"}`;
      }
    } catch (e) {
      throw e;
    }

  },
}