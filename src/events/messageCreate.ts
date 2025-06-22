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
    let text = message.content;
    if (message.attachments.size > 0) {
      const attachmentTypes: string[] = [];
      message.attachments.forEach((attachment) => {
        if (attachment?.contentType?.startsWith("audio/")) {
          attachmentTypes.push("音声ファイル");
          return;
        }
        if (attachment?.contentType?.startsWith("image/")) {
          attachmentTypes.push("画像ファイル");
          return;
        }
        if (attachment?.contentType?.startsWith("video/")) {
          attachmentTypes.push("動画ファイル");
          return;
        }
        if (attachment?.contentType?.startsWith("text/")) {
          attachmentTypes.push("テキストファイル");
          return;
        }
        attachmentTypes.push("不明なファイル");
      });
      if (attachmentTypes.length > 0) {
        text += `（${attachmentTypes.join("と")}が添付されました）`;
      }
    }
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEBOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEBOX_API_URL, headers);
    const query = await Query.getTalkQuery(text, 0);
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
