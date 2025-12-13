import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
  MessageFlags,
  PermissionFlagsBits,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";
import { ALStorage, loggingSystem } from "@/index";

export const handlingCommands = subCommandHandling("general/steeaming");
export const data = addSubCommand(
  "general/steeaming",
  new SlashCommandBuilder()
    .setName("steeaming")
    .setDescription("配信モードを管理します。")
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages)
);
export const guildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "general/steeaming" });

  const subcommand = (interaction.options as CommandInteractionOptionResolver).getSubcommand();
  const command = handlingCommands.get(subcommand);

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

    const errorMsg = (error as Error).toString();
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) {
      errorChannel.send({ embeds: [logEmbed] });
    }

    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    await interaction.reply({ embeds: [messageEmbed] });

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
};
