import {
  CommandInteraction,
  EmbedBuilder,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import { AudioPlayer, getVoiceConnection, VoiceConnectionReadyState } from "@discordjs/voice";
import { readTtsConnection } from "@/lib/dataUtils";

export const data = new SlashCommandSubcommandBuilder()
  .setName("skip")
  .setDescription("Skip the current audio.");

export const execute = async (interaction: CommandInteraction) => {
  const guildId = interaction.guild?.id;
  const channelId = interaction.channel?.id;
  if (!guildId || !channelId) {
    await interaction.reply({
      embeds: [
        new EmbedBuilder()
          .setTitle("Error - 情報不足")
          .setDescription("ギルドまたはチャンネル情報が取得できませんでした。")
          .setColor(interaction.client.config.color.error),
      ],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const voiceConnectionData = await readTtsConnection(guildId, channelId);
  if (!voiceConnectionData) {
    await interaction.reply({
      embeds: [
        new EmbedBuilder()
          .setTitle("Error - VC未接続")
          .setDescription("ボイスチャンネルに参加してないからスキップできません。")
          .setColor(interaction.client.config.color.error),
      ],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) {
    await interaction.reply({
      embeds: [
        new EmbedBuilder()
          .setTitle("Error - VC未接続")
          .setDescription("ボイスチャンネルに接続してないからスキップできません。")
          .setColor(interaction.client.config.color.error),
      ],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const player = (connection.state as VoiceConnectionReadyState).subscription?.player as
    | AudioPlayer
    | undefined;
  if (player) {
    player.stop(true);
  }

  await interaction.reply({
    embeds: [
      new EmbedBuilder()
        .setTitle("オーディオをスキップしました")
        .setDescription("今流れてたオーディオをスキップしました。")
        .setColor(interaction.client.config.color.success),
    ],
  });
};

export default {
  data,
  execute,
};
