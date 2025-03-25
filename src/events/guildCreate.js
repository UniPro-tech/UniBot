const {EmbedBuilder} = require("discord.js");
module.exports = {
  name: "guildCreate", // イベント名
  async execute(guild,client) {
    const channel = client.channels.cache.get(client.conf.logch.guildCreate);
    const log = new EmbedBuilder()
      .setTitle("GuildCreateLog")
      .setDescription(`Botが${guild.name}に参加しました。`)
      .setColor(client.conf.color.s)
      .setTimestamp()
      .setThumbnail(guild.iconURL({ dynamic: true }))
      .setFooter({ text: String(guild.id) });
    channel.send({ embeds: [log] });
  }
};
