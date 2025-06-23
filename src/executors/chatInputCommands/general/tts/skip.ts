import { CommandInteraction, MessageFlags, SlashCommandSubcommandBuilder } from "discord.js";
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
    await interaction.reply({
      content: "ボイスチャンネルに参加していません。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  const connection = getVoiceConnection(voiceConnectionData.guild);
  if (!connection) {
    await interaction.reply({
      content: "ボイスチャンネルに参加していません。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  const player = (connection.state as VoiceConnectionReadyState).subscription
    ?.player as AudioPlayer;
  if (player) {
    player.stop(true);
  }
  await interaction.reply({
    content: "現在のオーディオをスキップしました。",
    flags: [MessageFlags.Ephemeral],
  });
};

export default {
  data,
  execute,
};
