import {
  CommandInteraction,
  Embed,
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
  const voiceConnectionData = await readTtsConnection(
    interaction.guild?.id as string,
    interaction.channel?.id as string
  );
  if (!voiceConnectionData) {
    const embed = new EmbedBuilder()
      .setTitle("Error - VC未接続")
      .setDescription("ボイスチャンネルに参加していないため、オーディオをスキップできません。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }
  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) {
    const embed = new EmbedBuilder()
      .setTitle("Error - VC未接続")
      .setDescription("ボイスチャンネルに接続していないため、オーディオをスキップできません。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }
  const player = (connection.state as VoiceConnectionReadyState).subscription
    ?.player as AudioPlayer;
  if (player) {
    player.stop(true);
  }
  const embed = new EmbedBuilder()
    .setTitle("オーディオをスキップしました")
    .setDescription("現在のオーディオをスキップしました。")
    .setColor(interaction.client.config.color.success);
  await interaction.reply({ embeds: [embed] });
};

export default {
  data,
  execute,
};
