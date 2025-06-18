import { Client, EmbedBuilder, Guild, TextChannel } from "discord.js";
export const name = "guildCreate";
export const execute = async (guild: Guild, client: Client) => {
  const channel = client.channels.cache.get(client.config.logch.guildCreate);
  const log = new EmbedBuilder()
    .setTitle("GuildCreateLog")
    .setDescription(`Botが${guild.name}に参加しました。`)
    .setColor(client.config.color.success)
    .setTimestamp()
    .setThumbnail(guild.iconURL())
    .setFooter({ text: String(guild.id) });
  if (!channel || !(channel instanceof TextChannel)) {
    console.log("Error logChannel invalid");
    return;
  }
  channel.send({ embeds: [log] });
};

export default {
  name,
  execute,
};
