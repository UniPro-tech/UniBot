import {
  CommandInteraction,
  EmbedBuilder,
  GuildMember,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import { getVoiceConnection } from "@discordjs/voice";

export const data = new SlashCommandSubcommandBuilder()
  .setName("disconnect")
  .setDescription("Disconnect from the voice channel.");
export const execute = async (interaction: CommandInteraction) => {
  const channel = (interaction.member as GuildMember).voice.channel;
  if (!channel) {
    const embed = new EmbedBuilder()
      .setTitle("ボイスチャンネルに参加していません")
      .setDescription("ボイスチャンネルに参加していないため、接続を切断できません。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }
  const connection = await getVoiceConnection(channel.guild.id);
  if (!connection) {
    const embed = new EmbedBuilder()
      .setTitle("ボイスチャンネルに接続していません")
      .setDescription("ボイスチャンネルに接続していないため、切断できません。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }
  await connection.destroy();
  const embed = new EmbedBuilder()
    .setTitle("ボイスチャンネルから切断しました")
    .setDescription("ボイスチャンネルから切断しました。")
    .setColor(interaction.client.config.color.success);
  await interaction.reply({ embeds: [embed] });
};

export default {
  data,
  execute,
};
