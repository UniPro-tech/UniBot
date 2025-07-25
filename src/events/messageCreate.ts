import { listTtsDictionary, readTtsConnection, readTtsPreference } from "@/lib/dataUtils";
import {
  AudioPlayer,
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
  VoiceConnectionReadyState,
} from "@discordjs/voice";
import { Client, GuildChannel, Message } from "discord.js";
import { RPC, Generate, Query } from "voicevox.js";
import { Readable } from "stream";

export const name = "messageCreate";

const PLACEHOLDER = {
  CODEBLOCK: "__CODEBLOCK__",
  INLINECODE: "__INLINECODE__",
  LINK: "__LINK__",
};

const MAX_TEXT_LENGTH = 200;

function replacePlaceholders(text: string): string {
  return text
    .replace(/__CODEBLOCK_(\w+)__/g, (_, lang) => `、${lang}のコードブロック省略、`)
    .replace(/__CODEBLOCK__/g, "、コードブロック省略、")
    .replace(/__INLINECODE__/g, "、インラインコード省略、")
    .replace(/__LINK__/g, "、リンク省略、");
}

function escapeRegExp(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
}

async function replaceMentions(text: string, client: Client, message: Message): Promise<string> {
  // ユーザーメンション
  const userMentionRegex = /<@!?(\d+)>/g;
  const userMentions = [...text.matchAll(userMentionRegex)];
  for (const match of userMentions) {
    const userId = match[1];
    let replacement = "、ユーザー省略、";
    try {
      const user = await client.users.fetch(userId);
      replacement = `、${user.displayName}、`;
    } catch {}
    text = text.replace(match[0], replacement);
  }

  // ロールメンション
  const roleMentionRegex = /<@&(\d+)>/g;
  const roleMentions = [...text.matchAll(roleMentionRegex)];
  for (const match of roleMentions) {
    const roleId = match[1];
    let replacement = "、ロール省略、";
    try {
      const role = await message.guild?.roles.fetch(roleId);
      replacement = `、${role?.name ?? "ロール省略"}、`;
    } catch {}
    text = text.replace(match[0], replacement);
  }

  // チャンネルメンション
  const channelMentionRegex = /<#(\d+)>/g;
  const channelMentions = [...text.matchAll(channelMentionRegex)];
  for (const match of channelMentions) {
    const channelId = match[1];
    let replacement = "、チャンネル省略、";
    try {
      const channel = await client.channels.fetch(channelId);
      if (channel?.isTextBased()) replacement = `、${(channel as GuildChannel).name}、`;
    } catch {}
    text = text.replace(match[0], replacement);
  }

  return text;
}

function summarizeAttachments(message: Message): string {
  if (message.attachments.size === 0) return "";
  const typeCount: Record<string, number> = {};
  message.attachments.forEach((attachment) => {
    let type = "不明なファイル";
    if (attachment?.contentType?.startsWith("audio/")) type = "音声ファイル";
    else if (attachment?.contentType?.startsWith("image/")) type = "画像ファイル";
    else if (attachment?.contentType?.startsWith("video/")) type = "動画ファイル";
    else if (attachment?.contentType?.startsWith("text/")) type = "テキストファイル";
    typeCount[type] = (typeCount[type] || 0) + 1;
  });
  const typeStrings = Object.entries(typeCount).map(([type, count]) => `${count}個の${type}`);
  return typeStrings.length > 0 ? `（${typeStrings.join("と")}が添付されました）` : "";
}

function preprocessText(text: string): string {
  // コードブロック
  text = text.replace(/```(\w+)?\n?[\s\S]*?```/g, (m, lang) =>
    lang ? `__CODEBLOCK_${lang}__` : PLACEHOLDER.CODEBLOCK
  );
  // インラインコード
  text = text.replace(/`[^`]+`/g, PLACEHOLDER.INLINECODE);
  // リンク
  text = text.replace(/https?:\/\/\S+/g, PLACEHOLDER.LINK);
  return text;
}

function cutTextWithPlaceholders(text: string): string {
  const placeholderRegex = /__\w+?__/g;
  let cutIndex = 200;
  const matches = [...text.matchAll(placeholderRegex)];
  for (const match of matches) {
    const start = match.index!;
    const end = start + match[0].length;
    if (end > MAX_TEXT_LENGTH) {
      cutIndex = start;
      break;
    }
  }
  if (text.length > MAX_TEXT_LENGTH) {
    text = text.slice(0, cutIndex) + "（以下省略）";
  }
  return text;
}

function replaceSpecials(text: string): string {
  return text
    .replace(/<:.+?:\d+>/g, "、絵文字、")
    .replace(/<a:.+?:\d+>/g, "、アニメーション絵文字、")
    .replace(/<id:guide>/g, "、サーバーガイド、")
    .replace(/<id:browse>/g, "、チャンネル一覧、")
    .replace(/<id:customize>/g, "、チャンネルアンドロール、")
    .replace(/[\r\n]/g, "、");
}

async function applyDictionary(text: string, dict: any[]): Promise<string> {
  if (!Array.isArray(dict) || dict.length === 0) return text;
  for (const { word, definition, caseSensitive } of dict) {
    const regex = new RegExp(escapeRegExp(word), caseSensitive ? "g" : "gi");
    text = text.replace(regex, definition);
  }
  return text;
}

export const execute = async (message: Message, client: Client) => {
  if (message.author.bot || !message.guild || !message.channel.isTextBased()) return;
  if (!process.env.VOICEVOX_API_URL) {
    console.error("[ERROR] VOICEVOX_API_URL is not set.");
    return;
  }

  const voiceConnectionData = await readTtsConnection(message.guild.id, message.channel.id);
  if (!voiceConnectionData) return;

  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) return;

  if (message.flags.toArray().includes("SuppressNotifications")) return;

  if (["skip", "s"].includes(message.content)) {
    const player = (connection.state as VoiceConnectionReadyState).subscription
      ?.player as AudioPlayer;
    if (player) player.stop(true);
    return;
  }

  let text = preprocessText(message.content);
  text = cutTextWithPlaceholders(text);
  text = replacePlaceholders(text);

  text += summarizeAttachments(message);

  // メンション置換
  text = await replaceMentions(text, client, message);

  text = replaceSpecials(text);

  const dict = await listTtsDictionary(message.guild.id);
  text = await applyDictionary(text, dict);

  const styleId = ((await readTtsPreference(message.author.id, "speaker"))?.styleId as number) || 0;

  if (!RPC.rpc) {
    const headers = { Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}` };
    await RPC.connect(process.env.VOICEVOX_API_URL, headers);
  }

  const query = await Query.getTalkQuery(text, styleId);
  const audio = await Generate.generate(styleId, query);
  const audioStream = Readable.from(audio);
  const resource = createAudioResource(audioStream);

  let player: AudioPlayer | undefined = (connection.state as VoiceConnectionReadyState).subscription
    ?.player as AudioPlayer;

  if (player) {
    if (player.state.status === "playing") {
      await new Promise((resolve) => {
        player?.once("stateChange", (_, newState) => {
          if (newState.status === "idle") resolve(null);
        });
      });
    }
  } else {
    player = createAudioPlayer();
    connection.subscribe(player);
  }

  player.play(resource);
};

export default {
  name,
  execute,
};
