import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
  MessageFlags,
  PermissionFlagsBits,
  Collection,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";
import { ALStorage, loggingSystem } from "@/index";

// 初期化で循環参照を起こさないよう、subcommand の読み込みは遅延させる
export const handlingCommands = (() => {
  const col = new Collection<string, any>();
  setImmediate(() => {
    try {
      subCommandHandling("general/streaming", col);
    } catch (e) {
      // ログはトップレベルで使えない可能性があるため console を使う
      // 実行時のALStorageログは既存のロガーが利用される
      // eslint-disable-next-line no-console
      console.error("Failed to initialize handlingCommands for general/streaming", e);
    }
  });
  return col;
})();

export const data = (() => {
  const builder = new SlashCommandBuilder()
    .setName("streaming")
    .setDescription("配信モードを管理します。")
    .setDefaultMemberPermissions(PermissionFlagsBits.ManageMessages);
  setImmediate(() => {
    try {
      addSubCommand("general/streaming", builder);
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error("Failed to add subcommands for general/streaming", e);
    }
  });
  return builder;
})();
export const guildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "general/streaming" });

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
