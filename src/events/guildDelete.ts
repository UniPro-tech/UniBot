import { Client, EmbedBuilder, Guild, TextChannel } from "discord.js";

export const name = "guildDelete";
export const execute = async (guild: Guild, client: Client) => {
  const channel = client.channels.cache.get(client.config.logch.guildCreate);
  const log = new EmbedBuilder()
    .setTitle("GuildDeleteLog")
    .setDescription(`Botが${guild.name}にから退出しました。`)
    .setColor(client.config.color.s)
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
