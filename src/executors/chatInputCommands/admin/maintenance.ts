import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";
import { loggingSystem } from "@/index";

export const handlingCommands = subCommandHandling("admin/maintenance");
export const data = addSubCommand(
  "admin/maintenance",
  new SlashCommandBuilder().setName("maintenance").setDescription("メンテナンスモード")
);
export const guildOnly = true;
export const adminGuildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const logger = loggingSystem.getLogger({ function: "admin/maintenance" });
  const command = handlingCommands.get(
    (interaction.options as CommandInteractionOptionResolver).getSubcommand()
  );
  if (!command) {
    logger.error({ context: { command: interaction.commandName } }, "No command handler found");
    return;
  }
  try {
    await command.execute(interaction);
    logger.info({ context: { command: interaction.commandName } }, "Command executed successfully");
  } catch (error) {
    logger.error(
      { context: { command: interaction.commandName }, stack_trace: (error as Error).stack },
      "Command execution failed",
      error
    );
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const channel = await GetErrorChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません、エラーが発生しました...")
      .setDescription("```\n" + error + "\n```")
      .setColor(interaction.client.config.color.error)
      .setTimestamp();

    await interaction.reply({ embeds: [messageEmbed] });
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
  return;
};
