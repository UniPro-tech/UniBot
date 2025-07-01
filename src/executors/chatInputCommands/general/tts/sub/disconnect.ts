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
  const member = interaction.member as GuildMember;
  const channel = member.voice.channel;

  const errorEmbed = (desc: string) =>
    new EmbedBuilder()
      .setTitle("Error - VC未接続")
      .setDescription(desc)
      .setColor(interaction.client.config.color.error);

  if (!channel) {
    await interaction.reply({
      embeds: [errorEmbed("ボイスチャンネルに参加していないため、接続を切断できません。")],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const connection = getVoiceConnection(channel.guild.id);
  if (!connection) {
    await interaction.reply({
      embeds: [errorEmbed("ボイスチャンネルに接続していないため、切断できません。")],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  connection.destroy();

  const successEmbed = new EmbedBuilder()
    .setTitle("ボイスチャンネルから切断しました")
    .setDescription("ボイスチャンネルから切断しました。")
    .setColor(interaction.client.config.color.success);

  await interaction.reply({ embeds: [successEmbed] });
};

export default {
  data,
  execute,
};
