const { SlashCommandBuilder } = require("discord.js");
const Discord = require("discord.js");

module.exports = {
  guildOnly: false, // サーバー専用コマンドかどうか
  data: new SlashCommandBuilder() // スラッシュコマンド登録のため
    .setName("ping")
    .setDescription("Ping値を測定"),

  async execute(i, client) {
    const cmdPing = new Date() - i.createdAt;
    const embed = new Discord.EmbedBuilder()
      .setTitle("Ping")
      .setDescription("Pong!")
      .addFields([
        { name: 'WebSocket', value: ` ** ${client.ws.ping} ms ** `, inline: true },
        { name: 'コマンド受信', value: `** ${cmdPing} ms ** `, inline: true }])
      .setColor(client.conf.color.s)
      .setTimestamp();
    i.reply({ embeds: [embed] });
    return 'No data';
  },
}