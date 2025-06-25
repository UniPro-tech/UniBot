import { readTtsConnection } from "@/lib/dataUtils";
import {
  AudioPlayer,
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
  VoiceConnectionReadyState,
} from "@discordjs/voice";
import { Client, Message, MessageFlags } from "discord.js";
import { RPC, Generate, Query } from "voicevox.js";
import { Readable } from "stream";
export const name = "messageCreate";
export const execute = async (message: Message, client: Client) => {
  if (message.author.bot) return;
  if (!message.guild) return;
  const channel = message.channel;
  if (!channel.isTextBased()) return;
  if (process.env.VOICEVOX_API_URL === undefined) {
    console.error("[ERROR] VOICEVOX_API_URL is not set.");
    return;
  } else {
    const voiceConnectionData = await readTtsConnection(message.guild.id, channel.id);
    if (!voiceConnectionData) return;
    const connection = getVoiceConnection(voiceConnectionData.guild);
    if (!connection) return;
    if (message.flags.toArray().includes("SuppressNotifications")) return;
    if (message.content == "skip" || message.content == "s") {
      const player = (connection.state as VoiceConnectionReadyState).subscription
        ?.player as AudioPlayer;
      if (player) {
        player.stop(true);
      }
      return;
    }
    let text = "";
    if (message.content.length > 100) text = message.content.slice(0, 100) + "以下省略";
    else text = message.content;
    if (message.attachments.size > 0) {
      const typeCount: Record<string, number> = {};
      message.attachments.forEach((attachment) => {
        let type = "不明なファイル";
        if (attachment?.contentType?.startsWith("audio/")) {
          type = "音声ファイル";
        } else if (attachment?.contentType?.startsWith("image/")) {
          type = "画像ファイル";
        } else if (attachment?.contentType?.startsWith("video/")) {
          type = "動画ファイル";
        } else if (attachment?.contentType?.startsWith("text/")) {
          type = "テキストファイル";
        }
        typeCount[type] = (typeCount[type] || 0) + 1;
      });
      const typeStrings = Object.entries(typeCount).map(([type, count]) => `${count}個の${type}`);
      if (typeStrings.length > 0) {
        text += `（${typeStrings.join("と")}が添付されました）`;
      }
    }
    // まずMarkdown形式のリンクは、リンク部分だけ消して名前だけ残すよ！
    text = text.replace(/\[([^\]]+)\]\((https?:\/\/[^\s)]+)\)/g, "$1");
    // それ以外のhttp/httpsリンクは「リンク省略」にしとくね！
    text = text.replace(/https?:\/\/\S+/g, "、リンク省略、");
    text = text.replace(/[\r\n]/g, "、");
    if (!RPC.rpc) {
      const headers = {
        Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
      };
      await RPC.connect(process.env.VOICEVOX_API_URL, headers);
    }
    const query = await Query.getTalkQuery(text, 0);
    const audio = await Generate.generate(0, query);
    const audioStream = Readable.from(audio);
    const resource = createAudioResource(audioStream);
    let player: AudioPlayer | undefined = (connection.state as VoiceConnectionReadyState)
      .subscription?.player as AudioPlayer;
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
