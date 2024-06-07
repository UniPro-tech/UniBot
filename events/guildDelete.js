const {EmbedBuilder} = require("discord.js");
module.exports = {
  name: "guildDelete", // イベント名
  async execute(guild,client) {
    const channel = client.channels.cache.get(client.conf.logch.guildCreate);
    const log = new EmbedBuilder()
      .setTitle("GuildDeleteLog")
      .setDescription(`Botが${guild.name}にから退出しました。`)
      .setColor(client.conf.color.s)
      .setTimestamp()
      .setThumbnail(guild.iconURL({ dynamic: true }))
      .setFooter({ text: String(guild.id) });
    channel.send({ embeds: [log] });
  }
};
