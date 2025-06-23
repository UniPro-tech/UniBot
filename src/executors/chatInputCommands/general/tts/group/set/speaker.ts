import { CommandInteraction, MessageFlags, SlashCommandSubcommandBuilder } from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("speaker")
  .setDescription("Change the speaker of the TTS");
export const execute = async (interaction: CommandInteraction) => {
  // TODO: Implement the speaker change logic
  interaction.reply({
    content: "This feature is not implemented yet.",
    flags: [MessageFlags.Ephemeral],
  });
};

export default {
  data,
  execute,
};
