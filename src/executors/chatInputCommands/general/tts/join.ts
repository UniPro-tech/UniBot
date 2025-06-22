import { CommandInteraction, GuildMember, SlashCommandSubcommandBuilder } from "discord.js";
import { joinVoiceChannel } from "@discordjs/voice";
import { writeTtsConnection } from "@/lib/dataUtils";

export const data = new SlashCommandSubcommandBuilder()
  .setName("join")
  .setDescription("Connect to the voice channel.");
export const execute = async (interaction: CommandInteraction) => {
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
  execute,
};
