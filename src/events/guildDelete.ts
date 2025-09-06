import { Client, EmbedBuilder, Guild, TextChannel } from "discord.js";
import { ALStorage, loggingSystem } from "..";

export const name = "guildDelete";
export const execute = async (guild: Guild, client: Client) => {
  const ctx = { ...ALStorage.getStore(), context: { discord: { guild: guild.id } } };
  const logger = loggingSystem.getLogger({ ...ctx, function: "guildDelete" });
  const channel = client.channels.cache.get(client.config.logch.guildCreate);
  const log = new EmbedBuilder()
    .setTitle("GuildDeleteLog")
    .setDescription(`Botが${guild.name}にから退出しました。`)
    .setColor(client.config.color.success)
    .setTimestamp()
    .setThumbnail(guild.iconURL())
    .setFooter({ text: String(guild.id) });
  if (!channel || !(channel instanceof TextChannel)) {
    logger.error(
      { extra_context: { channelId: client.config.logch.guildCreate } },
      "GuildDeleteLog channel not found or is not a text channel."
    );
    return;
  }
  channel.send({ embeds: [log] });
};

export default {
  name,
  execute,
};
