import {
  ChatInputCommandInteraction,
  CommandInteractionOptionResolver,
  EmbedBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";

export const data = addSubCommand(
  "general/tts/group/set",
  new SlashCommandSubcommandGroupBuilder().setName("set").setDescription("Change TTS settings")
);
export const handlingCommand = subCommandHandling("general/tts/group/set");
export const execute = async (interaction: ChatInputCommandInteraction) => {
  const command = handlingCommand.get(
    (interaction.options as CommandInteractionOptionResolver).getSubcommand()
  );
  if (!command) {
    console.info(
      `[Not Found] Command: ${(
        interaction.options as CommandInteractionOptionResolver
      ).getSubcommand()}`
    );
    return;
  }
  try {
    await command.execute(interaction);
    console.info(
      `[Run] ${(interaction.options as CommandInteractionOptionResolver).getSubcommand()}`
    );
  } catch (error) {
    console.error(
      `[Error] ${(interaction.options as CommandInteractionOptionResolver).getSubcommand()}`,
      error
    );
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      const embed = new EmbedBuilder()
        .setTitle("TTS Command Error")
        .setDescription(`Error executing command: ${error}`)
        .setColor("Red");
      await logChannel.send({ embeds: [embed] });
    }
    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) {
      const embed = new EmbedBuilder()
        .setTitle("TTS Command Error")
        .setDescription(`Error executing command: ${error}`)
        .setColor("Red");
      await errorChannel.send({ embeds: [embed] });
    }
  }
  return;
};
