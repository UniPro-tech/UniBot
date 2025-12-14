import {
  ChannelType,
  CommandInteraction,
  EmbedBuilder,
  GuildMember,
  MessageFlags,
  SlashCommandSubcommandBuilder,
  TextChannel,
} from "discord.js";
import { getVoiceConnection, joinVoiceChannel, VoiceConnectionStatus } from "@discordjs/voice";
import { writeTtsConnection } from "@/lib/dataUtils";
import { TTSQueue } from "@/lib/ttsQueue";
import { ALStorage, loggingSystem } from "@/index";

export const data = new SlashCommandSubcommandBuilder()
  .setName("join")
  .setDescription("ボイスチャンネルに参加します。");

const createErrorEmbed = (title: string, description: string, color: number) =>
  new EmbedBuilder().setTitle(title).setDescription(description).setColor(color).setTimestamp();

const createSuccessEmbed = (voiceChannelId: string, textChannelId: string, color: number) =>
  new EmbedBuilder()
    .setTitle("TTSボイスチャンネル接続")
    .setDescription("ボイスチャンネルに接続しました。")
    .addFields([
      { name: "ボイスチャンネル名", value: `<#${voiceChannelId}>` },
      { name: "テキストチャンネル名", value: `<#${textChannelId}>` },
    ])
    .setColor(color);

export const execute = async (interaction: CommandInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "general/tts/join" });
  await interaction.deferReply();

  const guild = interaction.guild;
  const member = interaction.member as GuildMember;
  const voiceChannel = member.voice.channel;
  const config = interaction.client.config;

  if (!voiceChannel) {
    const embed = createErrorEmbed(
      "Error - TTSボイスチャンネル接続失敗",
      "ボイスチャンネルに参加していません。",
      config.color.error
    );
    await interaction.editReply({ embeds: [embed] });
    return;
  }

  if (getVoiceConnection(interaction.guild!.id)) {
    const embed = createErrorEmbed(
      "Error - 既に接続中",
      "既にボイスチャンネルに接続しています。",
      config.color.error
    );
    await interaction.editReply({ embeds: [embed] });
    return;
  }

  const connection = joinVoiceChannel({
    channelId: voiceChannel.id,
    guildId: voiceChannel.guild.id,
    adapterCreator: voiceChannel.guild.voiceAdapterCreator,
  });

  if (voiceChannel?.type === ChannelType.GuildStageVoice) {
    guild?.members.me?.voice.setSuppressed(false);
  }

  connection.once("ready", async () => {
    const textChannelId = (interaction.channel! as TextChannel).id;
    const embed = createSuccessEmbed(voiceChannel.id, textChannelId, config.color.success);
    await interaction.editReply({ embeds: [embed] });

    // TTS Queueに接続メッセージを追加（VoiceVox初期化待機付き）
    TTSQueue.enqueueConnectionMessage(voiceChannel.guild.id, voiceChannel.name).catch((error) => {
      logger.warn("Failed to enqueue connection message:", error as any);
    });

    let textChannel: string[] = [textChannelId];
    if (interaction.channel?.type !== ChannelType.GuildVoice) {
      textChannel.push(voiceChannel.id);
      voiceChannel.send({
        flags: [MessageFlags.SuppressNotifications],
        embeds: [
          new EmbedBuilder({
            title: "TTSボイスチャンネル接続",
            description: `このチャンネルに参加したためチャットにも接続しました。\n\nボイスチャンネル: <#${voiceChannel.id}>\nテキストチャンネル: <#${textChannelId}>`,
            color: config.color.success,
          }),
        ],
      });
    }
    writeTtsConnection(voiceChannel.guild.id, textChannel, voiceChannel.id);
  });

  connection.on("stateChange", (oldState, newState) => {
    logger.info(
      { extra_context: { guild: interaction.guild!.id, oldState, newState } },
      "Voice connection state changed"
    );
    if (newState.status === VoiceConnectionStatus.Disconnected) {
      const embed = createErrorEmbed(
        "Error - TTSボイスチャンネル切断",
        "ボイスチャンネルから切断されました。",
        config.color.error
      );
      (interaction.channel! as TextChannel).send({ embeds: [embed] });
      connection.destroy();
    }
  });
};

export default {
  data,
  execute,
};
