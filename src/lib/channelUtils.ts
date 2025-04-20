import { Client, TextChannel } from "discord.js";

export const GetLogChannel = async (client: Client) => {
  const channel = await client.channels.fetch(client.config.logch.command).catch((error) => null);
  if (!channel) {
    console.error(`[Not Found] Log Channel: ${client.config.logch.command}`);
    return null;
  }
  if (!(channel instanceof TextChannel)) {
    console.error(`[Type Error] Log Channel: ${client.config.logch.command}`);
    return null;
  }
  return channel;
};

export const GetErrorChannel = async (client: Client) => {
  const channel = await client.channels.fetch(client.config.logch.error).catch((error) => null);
  if (!channel) {
    console.error(`[Not Found] Log Channel: ${client.config.logch.command}`);
    return null;
  }
  if (!(channel instanceof TextChannel)) {
    console.error(`[Type Error] Log Channel: ${client.config.logch.command}`);
    return null;
  }
  return channel;
};

export default { GetLogChannel, GetErrorChannel };
