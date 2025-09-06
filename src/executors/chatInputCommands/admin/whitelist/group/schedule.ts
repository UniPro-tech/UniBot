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
  "admin/whitelist/group/schedule",
  new SlashCommandSubcommandGroupBuilder()
    .setName("schedule")
    .setDescription("スケジュールコマンドのホワイトリスト設定を変更")
);

export const handlingCommand = subCommandHandling("admin/whitelist/group/schedule");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const logger = loggingSystem.getLogger({ function: "admin/whitelist/group/schedule" });
  const options = interaction.options as CommandInteractionOptionResolver;
  const subcommand = options.getSubcommand();
  const command = handlingCommand.get(subcommand);

  if (!command) {
    logger.error(
      { context: { command: interaction.commandName, subcommand } },
      "No command handler found"
    );
    return;
  }

  try {
    await command.execute(interaction);
    logger.info(
      { context: { command: interaction.commandName, subcommand } },
      "Command executed successfully"
    );
  } catch (error) {
    logger.error(
      {
        context: { command: interaction.commandName, subcommand },
        stack_trace: (error as Error).stack,
      },
      "Command execution failed",
      error
    );

    const embed = new EmbedBuilder()
      .setTitle("Whitelist Command Error")
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
