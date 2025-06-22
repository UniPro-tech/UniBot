import { readTtsConnection } from "@/lib/dataUtils";
import { createAudioPlayer, createAudioResource, getVoiceConnection } from "@discordjs/voice";
import { Client, Message } from "discord.js";
import { RPC, Generate, Query } from "voicevox.js";
import { Readable } from "stream";
export const name = "messageCreate";
export const execute = async (message: Message, client: Client) => {
  if (message.author.bot) return;
  if (!message.guild) return;
  const channel = message.channel;
  if (!channel.isTextBased()) return;
  const voiceConnectionData = await readTtsConnection(message.guild.id, channel.id);
  if (!voiceConnectionData) return;
  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) return;
  if (process.env.VOICEBOX_API_URL === undefined) {
    console.error("VOICEBOX_API_URL is not set.");
    return;
  } else {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEBOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEBOX_API_URL, headers);
    const query = await Query.getTalkQuery(message.content, 0);
    const audio = await Generate.generate(0, query);
    const player = createAudioPlayer();
    connection.subscribe(player);
    const audioStream = Readable.from(audio);
    const resource = createAudioResource(audioStream);
    player.play(resource);
  }
};

export default {
  name,
  execute,
};
