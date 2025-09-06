import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, MessageContextMenuCommandInteraction } from "discord.js";
import config from "@/config";
import { loggingSystem } from "@/index";

const MessageContextMenuCommandExecute = async (
  interaction: MessageContextMenuCommandInteraction
) => {
  const logger = loggingSystem.getLogger({ function: "MessageContextMenuCommandExecute" });
  logger.info(
    { extra_context: { commandName: interaction.commandName } },
    "MessageContextMenuCommand execution started"
  );

  const command = interaction.client.interactionExecutorsCollections.messageContextMenuCommands.get(
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
    logger.info(
      { extra_context: { commandName: interaction.commandName } },
      "MessageContextMenuCommand executed successfully"
    );

    const logEmbed = new EmbedBuilder()
      .setTitle("コマンド実行ログ")
      .setDescription(`${interaction.user} がコマンドを実行したよ！`)
      .setColor(config.color.success)
      .setTimestamp()
      .setThumbnail(interaction.user.displayAvatarURL())
      .addFields([
        {
          name: "コマンド",
          value: `\`\`\`\n/${interaction.commandName}\n\`\`\``,
        },
        {
          name: "実行サーバー",
          value:
            "```\n" +
            (interaction.inGuild()
              ? `${interaction.guild?.name} (${interaction.guild?.id})`
              : "DM") +
            "\n```",
        },
        {
          name: "実行ユーザー",
          value: "```\n" + `${interaction.user.tag}(${interaction.user.id})` + "\n```",
        },
      ])
      .setFooter({ text: `${interaction.id}` });

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) await logChannel.send({ embeds: [logEmbed] });
  } catch (error) {
    logger.error(
      {
        extra_context: { commandName: interaction.commandName },
        stack_trace: (error as Error).stack,
      },
      "Command execution failed",
      error
    );

    const errorEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) await errorChannel.send({ embeds: [errorEmbed] });

    const userEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    await interaction.reply({ embeds: [userEmbed] }).catch(() => {});
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) await logChannel.send({ embeds: [userEmbed] });
  }
};

export default MessageContextMenuCommandExecute;
