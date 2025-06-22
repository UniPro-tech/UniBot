import { readTtsConnection } from "@/lib/dataUtils";
import { getVoiceConnection } from "@discordjs/voice";
import { EmbedBuilder, TextChannel, VoiceState } from "discord.js";
export const name = "voiceStateUpdate";
export const execute = async (oldState: VoiceState) => {
  const channel = oldState.channel;
  if (!channel) return;
  if (channel.members.size > 1) return;
  const connectionData = await readTtsConnection(oldState.guild.id, channel.id);
  if (!connectionData) return;
  const connection = getVoiceConnection(channel.id);
  if (!connection || connection.state.status == "destroyed") return;
  connection.destroy();
  const textChannel = oldState.guild.channels.cache.get(
    connectionData.textChannel[0] as string
  ) as TextChannel;
  textChannel.send("ボイスチャンネルから誰もいなくなったので、接続を切断しました。");
  const logEmbed = new EmbedBuilder()
    .setTitle("ボイスチャンネル切断ログ")
    .setDescription(`<#${oldState.channel.id}> が無人になったため切断しました。`)
    .setColor(oldState.client.config.color.success)
    .setTimestamp();
  const logChannel = oldState.client.channels.cache.get(
    oldState.client.config.logch.command
  ) as TextChannel;
  if (!logChannel) {
    console.error("[Error] LogChannel invalid");
    return;
  }
  logChannel.send({ embeds: [logEmbed] });
};

export default {
  name,
  execute,
};
