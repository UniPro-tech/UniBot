import { Client, EmbedBuilder, Guild, TextChannel } from "discord.js";
import { loggingSystem } from "..";
export const name = "guildCreate";
export const execute = async (guild: Guild, client: Client) => {
  const logger = loggingSystem.getLogger({ function: "guildCreate" });
  const channel = client.channels.cache.get(client.config.logch.guildCreate);
  const log = new EmbedBuilder()
    .setTitle("GuildCreateLog")
    .setDescription(`Botが${guild.name}に参加しました。`)
    .setColor(client.config.color.success)
    .setTimestamp()
    .setThumbnail(guild.iconURL())
    .setFooter({ text: String(guild.id) });
  if (!channel || !(channel instanceof TextChannel)) {
    logger.error(
      { context: { channelId: client.config.logch.guildCreate } },
      "GuildCreateLog channel not found or is not a text channel."
    );
    return;
  }
  channel.send({ embeds: [log] });
};

export default {
  name,
  execute,
};
