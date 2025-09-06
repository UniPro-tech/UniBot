import {
  ChatInputCommandInteraction,
  CommandInteractionOptionResolver,
  EmbedBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { loggingSystem } from "@/index";

export const data = addSubCommand(
  "general/tts/group/set",
  new SlashCommandSubcommandGroupBuilder().setName("set").setDescription("TTSの設定を変更")
);

export const handlingCommand = subCommandHandling("general/tts/group/set");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const logger = loggingSystem.getLogger({ function: "general/tts/set" });
  const options = interaction.options as CommandInteractionOptionResolver;
  const subcommand = options.getSubcommand();
  const command = handlingCommand.get(subcommand);

  if (!command) {
    logger.error(
      { extra_context: { command: interaction.commandName } },
      "No command handler found"
    );
    return;
  }

  try {
    await command.execute(interaction);
    logger.info(
      { extra_context: { command: interaction.commandName } },
      "Command executed successfully"
    );
  } catch (error) {
    logger.error(
      { extra_context: { command: interaction.commandName }, stack_trace: (error as Error).stack },
      "Command execution failed",
      error
    );

    const embed = new EmbedBuilder()
      .setTitle("TTS Command Error")
      .setDescription(`Error executing command: ${error}`)
      .setColor("Red");

    const [logChannel, errorChannel] = await Promise.all([
      GetLogChannel(interaction.client),
      GetErrorChannel(interaction.client),
    ]);

    if (logChannel) await logChannel.send({ embeds: [embed] });
    if (errorChannel) await errorChannel.send({ embeds: [embed] });
  }
};
