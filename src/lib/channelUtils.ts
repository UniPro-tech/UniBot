import { Client, TextChannel } from "discord.js";
import { loggingSystem } from "..";

export const GetLogChannel = async (client: Client) => {
  const logger = loggingSystem.getLogger({ function: "GetLogChannel" });
  const channel = await client.channels.fetch(client.config.logch.command).catch((error) => null);
  if (!channel) {
    logger.error(
      { extra_context: { channelId: client.config.logch.command } },
      `Log Channel not found`
    );
    return null;
  }
  if (!(channel instanceof TextChannel)) {
    logger.error(
      { extra_context: { channelId: client.config.logch.command } },
      `Log Channel type error`
    );
    return null;
  }
  return channel;
};

export const GetErrorChannel = async (client: Client) => {
  const logger = loggingSystem.getLogger({ function: "GetErrorChannel" });
  const channel = await client.channels.fetch(client.config.logch.error).catch((error) => null);
  if (!channel) {
    logger.error(
      { extra_context: { channelId: client.config.logch.error } },
      `Error Channel not found`
    );
    return null;
  }
  if (!(channel instanceof TextChannel)) {
    logger.error(
      { extra_context: { channelId: client.config.logch.error } },
      `Error Channel type error`
    );
    return null;
  }
  return channel;
};

export default { GetLogChannel, GetErrorChannel };
