import { readTtsConnection } from "@/lib/dataUtils";
import {
  AudioPlayer,
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
  VoiceConnectionReadyState,
} from "@discordjs/voice";
import { Client, GuildChannel, Message, Role } from "discord.js";
import { RPC, Generate, Query } from "voicevox.js";
import { Readable } from "stream";

export const name = "messageCreate";
export const execute = async (message: Message, client: Client) => {
  if (message.author.bot) return;
  if (!message.guild) return;
  if (!message.channel.isTextBased()) return;

  if (!process.env.VOICEVOX_API_URL) {
    console.error("[ERROR] VOICEVOX_API_URL is not set.");
    return;
  }

  const voiceConnectionData = await readTtsConnection(message.guild.id, message.channel.id);
  if (!voiceConnectionData) return;

  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) return;

  if (message.flags.toArray().includes("SuppressNotifications")) return;

  if (message.content === "skip" || message.content === "s") {
    const player = (connection.state as VoiceConnectionReadyState).subscription
      ?.player as AudioPlayer;
    if (player) player.stop(true);
    return;
  }

  // 元のテキスト
  let text = message.content;

  // 1. 先にコードブロック、インラインコード、リンクをプレースホルダーに置換
  text = text.replace(/```(\w+)?\n?[\s\S]*?```/g, (m, lang) =>
    lang ? `__CODEBLOCK_${lang}__` : "__CODEBLOCK__"
  );
  text = text.replace(/`[^`]+`/g, "__INLINECODE__");
  text = text.replace(/https?:\/\/\S+/g, "__LINK__");

  // 2. 文字数制限（プレースホルダー含む）
  // プレースホルダーの位置を全部取る
  const placeholderRegex = /__\w+?__/g;
  let cutIndex = 100;

  const matches = [...text.matchAll(placeholderRegex)];
  for (const match of matches) {
    const start = match.index!;
    const end = start + match[0].length;
    if (end > 200) {
      // プレースホルダーの途中で切れるなら、切る位置をプレースホルダーの終わりにずらす
      cutIndex = start; // 途中で切るより直前で切るほうが安全
      break;
    }
  }

  // 文字数制限カット
  if (text.length > 200) {
    text = text.slice(0, cutIndex) + "（以下省略）";
  }

  // 3. プレースホルダーを読み上げ用テキストに変換
  text = text.replace(/__CODEBLOCK_(\w+)__/g, (m, lang) => `、${lang}のコードブロック省略、`);
  text = text.replace(/__CODEBLOCK__/g, "、コードブロック省略、");
  text = text.replace(/__INLINECODE__/g, "、インラインコード省略、");
  text = text.replace(/__LINK__/g, "、リンク省略、");

  // 4. 添付ファイルの種類カウント
  if (message.attachments.size > 0) {
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
    if (typeStrings.length > 0) {
      text += `（${typeStrings.join("と")}が添付されました）`;
    }
  }

  // 5. メンション類の置換

  // ユーザーメンション
  const mentionRegex = /<@!?(\d+)>/g;
  const mentionMatches = [...text.matchAll(mentionRegex)];
  for (const match of mentionMatches) {
    const userId = match[1];
    let username = "ユーザー省略";
    try {
      const user = await client.users.fetch(userId);
      username = user.username;
    } catch {}
    text = text.replace(match[0], `、${username}、`);
  }

  // ロールメンション
  const roleMentionRegex = /<@&(\d+)>/g;
  const roleMentionMatches = [...text.matchAll(roleMentionRegex)];
  for (const match of roleMentionMatches) {
    const roleId = match[1];
    let roleName = "ロール省略";
    try {
      const role = await message.guild.roles.fetch(roleId);
      if (role) roleName = role.name;
    } catch {}
    text = text.replace(match[0], `、${roleName}、`);
  }

  // チャンネルメンション
  const channelMentionRegex = /<#[0-9]+>/g;
  const channelMentionMatches = [...text.matchAll(channelMentionRegex)];
  for (const match of channelMentionMatches) {
    const channelId = match[0].slice(2, -1);
    let channelName = "チャンネル省略";
    try {
      const channel = await client.channels.fetch(channelId);
      if (channel?.isTextBased()) channelName = (channel as GuildChannel).name;
    } catch {}
    text = text.replace(match[0], `、${channelName}、`);
  }

  // 6. 絵文字・サーバーガイド置換
  text = text.replace(/<:.+?:\d+>/g, "、絵文字、");
  text = text.replace(/<a:.+?:\d+>/g, "、アニメーション絵文字、");
  text = text.replace(/<id:guide>/g, "、サーバーガイド、");
  text = text.replace(/<id:browse>/g, "、チャンネル一覧、");
  text = text.replace(/<id:customize>/g, "、チャンネルアンドロール、");

  // 7. 改行をカンマに変換
  text = text.replace(/[\r\n]/g, "、");

  // VOICEVOXの接続と音声再生
  if (!RPC.rpc) {
    const headers = { Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}` };
    await RPC.connect(process.env.VOICEVOX_API_URL, headers);
  }

  const query = await Query.getTalkQuery(text, 0);
  const audio = await Generate.generate(0, query);
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
