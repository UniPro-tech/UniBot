import {
  CommandInteraction,
  GuildMember,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";
import { joinVoiceChannel } from "@discordjs/voice";
import { writeTtsConnection } from "@/lib/dataUtils";

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
  const voiceChannel = (interaction.member as GuildMember).voice.channel;
  if (!voiceChannel) {
    await interaction.followUp("ボイスチャンネルに参加していません。");
    return;
  }
  const connection = await joinVoiceChannel({
    channelId: voiceChannel.id,
    guildId: voiceChannel.guild.id,
    adapterCreator: voiceChannel.guild.voiceAdapterCreator,
  });
  connection.once("ready", async () => {
    console.info("Connected to voice channel");
    await interaction.editReply("Connection OK");
    writeTtsConnection(voiceChannel.guild.id, [interaction.channel?.id as string], voiceChannel.id);
  });
};

export default {
  data,
  adminGuildOnly,
  execute,
};
