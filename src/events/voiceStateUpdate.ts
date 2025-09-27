import { readTtsConnection } from "@/lib/dataUtils";
import { TTSQueue } from "@/lib/ttsQueue";
import { getVoiceConnection } from "@discordjs/voice";
import { EmbedBuilder, TextChannel, VoiceBasedChannel, VoiceState } from "discord.js";
import { ALStorage, loggingSystem } from "..";

export const name = "voiceStateUpdate";

const sendEmbed = async (
  channel: TextChannel,
  title: string,
  description: string,
  color: number
) => {
  const embed = new EmbedBuilder()
    .setTitle(title)
    .setDescription(description)
    .setColor(color)
    .setTimestamp();
  await channel.send({ embeds: [embed] });
};

const getCurrentChannel = (
  oldState: VoiceState,
  newState: VoiceState
): VoiceBasedChannel | null | undefined => {
  if (newState.channel) {
    return newState.guild.members.cache.get(newState.client.user?.id!)?.voice?.channel;
  }
  if (oldState.channel) {
    return oldState.guild.members.cache.get(oldState.client.user?.id!)?.voice?.channel;
  }
  return null;
};

const handleDisconnect = async (
  oldState: VoiceState,
  currentChannel: VoiceBasedChannel | null | undefined
) => {
  const ctx = {
    ...ALStorage.getStore(),
    user_id: oldState.member?.user.id,
    context: { discord: { guild: oldState.guild.id, channel: oldState.channel?.id } },
  };
  const logger = loggingSystem.getLogger({ ...ctx, function: "handleDisconnect" });
  if (!currentChannel || oldState.channel?.id !== currentChannel.id) return;
  const connectionData = await readTtsConnection(oldState.guild.id, undefined, currentChannel.id);
  if (!connectionData) return;
  const connection = getVoiceConnection(oldState.guild.id);
  if (!connection || connection.state.status === "destroyed") return;
  connection.destroy();

  // TTSQueueインスタンスを削除
  TTSQueue.removeInstance(oldState.guild.id);

  const textChannel = oldState.guild.channels.cache.get(
    connectionData.textChannel[0] as string
  ) as TextChannel;
  await sendEmbed(
    textChannel,
    "ボイスチャンネル切断",
    `<#${oldState.channel!.id}> が無人になったため切断しました。`,
    oldState.client.config.color.success
  );

  const logChannel = oldState.client.channels.cache.get(
    oldState.client.config.logch.command
  ) as TextChannel;
  if (!logChannel) {
    logger.error(
      { extra_context: { channel: oldState.client.config.logch.command } },
      "Log channel not found."
    );
    return;
  }
  await sendEmbed(
    logChannel,
    "ボイスチャンネル切断ログ",
    `<#${oldState.channel!.id}> が無人になったため切断しました。`,
    oldState.client.config.color.success
  );
};

const getVoiceEventType = (oldState: VoiceState, newState: VoiceState) => {
  if (newState.channel && oldState.channel) return "switch";
  if (newState.channel) return "join";
  return "leave";
};

const getVoiceEventText = (type: string, oldState: VoiceState, newState: VoiceState) => {
  switch (type) {
    case "switch":
      return `${newState.member?.displayName} が ${oldState.channel?.name} から ${newState.channel?.name} に切り替えました。`;
    case "join":
      return `${newState.member?.displayName} が ${newState.channel?.name} に参加しました。`;
    case "leave":
      return `${oldState.member?.displayName} が ${oldState.channel?.name} から退出しました。`;
    default:
      return "";
  }
};

const speak = async (guildId: string, text: string) => {
  const ttsQueue = TTSQueue.getInstance(guildId);
  // ボイスチャンネル通知は高優先度でキューに追加
  ttsQueue.enqueue(text, 0, 0); // styleId: 0, priority: 0 (高優先度)
};

export const execute = async (oldState: VoiceState, newState: VoiceState) => {
  const oldChannel = oldState.channel;
  const newChannel = newState.channel;
  const currentChannel = getCurrentChannel(oldState, newState);
  if (!currentChannel) return;
  if (!(oldChannel?.id === currentChannel.id || newChannel?.id === currentChannel.id)) return;

  if (oldChannel && oldChannel.members.size === 1) {
    await handleDisconnect(oldState, currentChannel);
    return;
  }

  if (
    newState.member?.user.id === newState.client.user?.id ||
    oldState.member?.user.id === oldState.client.user?.id
  )
    return;
  if (newChannel?.id === oldChannel?.id) return;

  const type = getVoiceEventType(oldState, newState);
  const text = getVoiceEventText(type, oldState, newState);

  const connection = getVoiceConnection(oldState.guild.id || newState.guild.id);
  if (connection && connection.state.status !== "destroyed") {
    await speak(oldState.guild.id || newState.guild.id, text);
  }
};

export default {
  name,
  execute,
};
