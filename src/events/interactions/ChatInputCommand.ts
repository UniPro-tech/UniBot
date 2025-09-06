import config from "@/config";
import { loggingSystem } from "@/index";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { writeChatInputCommandInteractionLog } from "@/lib/logger";
import { ChatInputCommandInteraction, EmbedBuilder } from "discord.js";

const ChatInputCommandExecute = async (interaction: ChatInputCommandInteraction) => {
  const logger = loggingSystem.getLogger({ function: "ChatInputCommandExecute" });
  logger.info(
    { extra_context: { commandName: interaction.commandName } },
    "ChatInputCommand execution started"
  );

  const command = interaction.client.interactionExecutorsCollections.chatInputCommands.get(
    interaction.commandName
  );
  if (!command) {
    logger.error(
      { extra_context: { commandName: interaction.commandName } },
      "No command handler found for this command"
    );
    return;
  }

  if (!interaction.inGuild() && command.guildOnly) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはDMでは実行できません。")
      .setColor(interaction.client.config.color.error);

    await interaction.reply({ embeds: [embed] });
    logger.info(
      { extra_context: { commandName: interaction.commandName } },
      "Blocked command execution in DM"
    );
    return;
  }

  try {
    await command.execute(interaction);
    await writeChatInputCommandInteractionLog(interaction);
  } catch (error) {
    logger.error(
      {
        extra_context: { commandName: interaction.commandName },
        stack_trace: (error as Error).stack,
      },
      "Command execution failed",
      error
    );

    const errorMsg = (error as any).toString();
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) await errorChannel.send({ embeds: [logEmbed] });

    const userEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    if (!interaction.replied && !interaction.deferred) {
      await interaction.reply({ embeds: [userEmbed] });
    } else {
      await interaction.followUp({ embeds: [userEmbed] });
    }

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) await logChannel.send({ embeds: [userEmbed] });
  }
};

export default ChatInputCommandExecute;
