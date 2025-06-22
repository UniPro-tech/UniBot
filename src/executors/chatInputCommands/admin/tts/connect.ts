import {
  CommandInteraction,
  GuildMember,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";
import { joinVoiceChannel } from "@discordjs/voice";

export const data = new SlashCommandSubcommandBuilder()
  .setName("connect")
  .setDescription("Connect to the voice channel.");
export const adminGuildOnly = true;
export const execute = async (interaction: CommandInteraction) => {
  if (!(interaction.member?.roles as GuildMemberRoleManager).cache.has(config.adminRoleId)) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("Now connecting...");
  const channel = (interaction.member as GuildMember).voice.channel;
  if (!channel) {
    await interaction.followUp("ボイスチャンネルに参加していません。");
    return;
  }
  const connection = await joinVoiceChannel({
    channelId: channel.id,
    guildId: channel.guild.id,
    adapterCreator: channel.guild.voiceAdapterCreator,
  });
  await interaction.editReply("Connection OK");
};

export default {
  data,
  adminGuildOnly,
  execute,
};
