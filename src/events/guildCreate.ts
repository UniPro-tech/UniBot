import { Client, EmbedBuilder, Guild, TextChannel } from "discord.js";
import { ALStorage, loggingSystem } from "..";
export const name = "guildCreate";
export const execute = async (guild: Guild, client: Client) => {
  const ctx = { ...ALStorage.getStore(), context: { discord: { guild: guild.id } } };
  const logger = loggingSystem.getLogger({ ...ctx, function: "guildCreate" });
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
      { extra_context: { channelId: client.config.logch.guildCreate } },
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
