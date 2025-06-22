import {
  CommandInteraction,
  GuildMember,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";
import { getVoiceConnection, joinVoiceChannel } from "@discordjs/voice";

export const data = new SlashCommandSubcommandBuilder()
  .setName("disconnect")
  .setDescription("Disconnect from the voice channel.");
export const adminGuildOnly = true;
export const execute = async (interaction: CommandInteraction) => {
  if (!(interaction.member?.roles as GuildMemberRoleManager).cache.has(config.adminRoleId)) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("Now disconnecting...");
  const channel = (interaction.member as GuildMember).voice.channel;
  if (!channel) {
    await interaction.followUp("ボイスチャンネルに参加していません。");
    return;
  }
  const connection = await getVoiceConnection(channel.guild.id);
  if (!connection) {
    await interaction.followUp("ボイスチャンネルに接続していません。");
    return;
  }
  connection.destroy();
  await interaction.editReply("Connection disconnected successfully.");
};

export default {
  data,
  adminGuildOnly,
  execute,
};
