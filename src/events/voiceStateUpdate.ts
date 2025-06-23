import { readTtsConnection } from "@/lib/dataUtils";
import { createAudioPlayer, createAudioResource, getVoiceConnection } from "@discordjs/voice";
import { ChannelType, EmbedBuilder, TextChannel, VoiceState } from "discord.js";
import { Readable } from "stream";
import { RPC, Query, Generate } from "voicevox.js";
export const name = "voiceStateUpdate";
export const execute = async (oldState: VoiceState, newState: VoiceState) => {
  const channel = oldState.channel;
  if (channel && channel.members.size === 1) {
    if (channel.type !== ChannelType.GuildVoice) return;
    const connectionData = await readTtsConnection(oldState.guild.id, undefined, channel.id);
    if (!connectionData) return;
    const connection = getVoiceConnection(oldState.guild.id);
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
    return;
  }
  if (newState.member?.user.bot || oldState.member?.user.bot) return;
  const type = newState.channel ? (oldState.channel ? "switch" : "join") : "leave";
  const text =
    type === "switch"
      ? `${newState.member?.displayName} がボイスチャンネルを ${oldState.channel?.name} から ${newState.channel?.name} に切り替えました。`
      : type === "join"
      ? `${newState.member?.displayName} がボイスチャンネル ${newState.channel?.name} に参加しました。`
      : `${oldState.member?.displayName} がボイスチャンネル ${oldState.channel?.name} から退出しました。`;
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
      while (player.state.status === "playing") {
        await new Promise((resolve) => setTimeout(resolve, 50));
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
