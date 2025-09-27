import {
  listTtsDictionary,
  readTtsConnection,
  readTtsPreference,
  ServerDataManager,
} from "@/lib/dataUtils";
import { TTSQueue } from "@/lib/ttsQueue";
import { getVoiceConnection } from "@discordjs/voice";
import { Client, EmbedBuilder, GuildChannel, Message, PartialGroupDMChannel } from "discord.js";
import { ALStorage, loggingSystem } from "..";
import config from "@/config";

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

async function voicevoxSynthesis(message: Message, client: Client): Promise<void> {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "voicevoxSynthesis" });
  if (message.author.bot || !message.guild || !message.channel.isTextBased()) return;
  if (!process.env.VOICEVOX_API_URL) {
    logger.warn("VOICEVOX_API_URL is not set.");
    return;
  }

  const voiceConnectionData = await readTtsConnection(message.guild.id, message.channel.id);
  if (!voiceConnectionData) return;

  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) return;

  if (message.flags.toArray().includes("SuppressNotifications")) return;

  // skipコマンドの処理
  if (["skip", "s"].includes(message.content)) {
    const ttsQueue = TTSQueue.getInstance(message.guild.id);
    const skipped = ttsQueue.skip();
    if (skipped) {
      logger.debug(`TTS skipped by user ${message.author.id} in guild ${message.guild.id}`);
    }
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

  // TTSQueueに追加
  const ttsQueue = TTSQueue.getInstance(message.guild.id);
  ttsQueue.enqueue(text, styleId, 1);
  logger.debug(
    `Enqueued TTS for user ${message.author.id} in guild ${message.guild.id}: ${text} (styleId: ${styleId})`
  );
}

const resendPinnedMessage = async (message: Message, client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "resendPinnedMessage" });
  if (
    !message.guild ||
    !message.channel.isTextBased() ||
    message.channel instanceof PartialGroupDMChannel ||
    (message.author.id === client.user?.id &&
      (message.embeds.length === 0 || message.embeds[0].footer?.text.includes("Pinned Message")))
  )
    return;

  const channelId = message.channel.id;
  const dataManager = new ServerDataManager(message.guild.id);
  const pinnedMessageConfig = await dataManager.getConfig("pinnedMessage", channelId);
  if (!pinnedMessageConfig || !pinnedMessageConfig.message) return;
  try {
    const oldMessage = await message.channel.messages.fetch(pinnedMessageConfig.latestMessageId);
    if (oldMessage) await oldMessage.delete();
    const embed = new EmbedBuilder()
      .setDescription(pinnedMessageConfig.message)
      .setColor(config.color.success)
      .setFooter({ text: "Pinned Message" });
    const newSendedMessage = await message.channel.send({ embeds: [embed] });
    dataManager.setConfig(
      "pinnedMessage",
      { message: pinnedMessageConfig.message, latestMessageId: newSendedMessage.id },
      channelId
    );
  } catch (error) {
    logger.error("Failed to resend pinned message:", error as any);
  }
};

export const execute = async (message: Message, client: Client) => {
  voicevoxSynthesis(message, client).catch((error) => {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "voicevoxSynthesis" });
    logger.error(error);
  });
  resendPinnedMessage(message, client).catch((error) => {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "resendPinnedMessage" });
    logger.error("Error in resendPinnedMessage:", error);
  });
};

export default {
  name,
  execute,
};
