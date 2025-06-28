import { readTtsConnection } from "@/lib/dataUtils";
import {
  AudioPlayer,
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
} from "@discordjs/voice";
import { EmbedBuilder, TextChannel, VoiceState } from "discord.js";
import { Readable } from "stream";
import { RPC, Query, Generate } from "voicevox.js";
export const name = "voiceStateUpdate";
export const execute = async (oldState: VoiceState, newState: VoiceState) => {
  const oldChannel = oldState.channel;
  const newChannel = newState.channel;
  const currentChannel = newChannel
    ? newState.guild.members.cache.get(newState.client.user?.id)?.voice?.channel
    : oldChannel
    ? oldState.guild.members.cache.get(oldState.client.user?.id)?.voice?.channel
    : null;
  if (!currentChannel) return;
  if (!(oldChannel?.id === currentChannel.id || newChannel?.id === currentChannel.id)) return;
  if (oldState.channel && oldState.channel.members.size === 1) {
    const connectionData = await readTtsConnection(oldState.guild.id, undefined, currentChannel.id);
    if (!connectionData) return;
    const connection = getVoiceConnection(oldState.guild.id);
    if (!connection || connection.state.status == "destroyed") return;
    connection.destroy();
    const textChannel = oldState.guild.channels.cache.get(
      connectionData.textChannel[0] as string
    ) as TextChannel;
    const embed = new EmbedBuilder()
      .setTitle("ボイスチャンネル切断")
      .setDescription(`<#${oldState.channel.id}> が無人になったため切断しました。`)
      .setColor(oldState.client.config.color.success)
      .setTimestamp();
    textChannel.send({ embeds: [embed] });
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
    return;
  }
  if (
    newState.member?.user.id === newState.client.user?.id ||
    oldState.member?.user.id === oldState.client.user?.id
  )
    return;
  if (newState.channel?.id === oldState.channel?.id) return;

  const type = newState.channel ? (oldState.channel ? "switch" : "join") : "leave";
  const text =
    type === "switch"
      ? `${newState.member?.displayName} が ${oldState.channel?.name} から ${newState.channel?.name} に切り替えました。`
      : type === "join"
      ? `${newState.member?.displayName} が ${newState.channel?.name} に参加しました。`
      : `${oldState.member?.displayName} が ${oldState.channel?.name} から退出しました。`;
  const connection = getVoiceConnection(oldState.guild.id || newState.guild.id);
  if (connection && connection.state.status !== "destroyed") {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    };
    if (!RPC.rpc) await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
    const query = await Query.getTalkQuery(text, 0);
    const audio = await Generate.generate(0, query);
    const audioStream = Readable.from(audio);
    const resource = createAudioResource(audioStream);
    let player = connection.state.subscription?.player;
    if (player) {
      if (player.state.status === "playing") {
        await new Promise((resolve) => {
          (player as AudioPlayer).once("stateChange", (oldState, newState) => {
            if (newState.status === "idle") {
              resolve(null);
            }
          });
        });
      }
    } else {
      player = createAudioPlayer();
      connection.subscribe(player);
    }
    player.play(resource);
  }
};

export default {
  name,
  execute,
};
