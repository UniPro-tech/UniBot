import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
  TextChannel,
  MessageFlags,
  InteractionContextType,
} from "discord.js";
import { addSubCommand, addSubCommandGroup, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";
import { ALStorage, loggingSystem } from "@/index";

export const handlingCommands = subCommandHandling(
  "general/tts/group",
  subCommandHandling("general/tts/sub")
);

export const data = addSubCommandGroup(
  "general/tts/group",
  addSubCommand(
    "general/tts/sub",
    new SlashCommandBuilder().setName("tts").setDescription("テキスト読み上げを管理します。")
  ) as SlashCommandBuilder
).setContexts(InteractionContextType.Guild);

export const adminGuildOnly = true;

const replyWithError = async (
  interaction: ChatInputCommandInteraction,
  title: string,
  description: string
) => {
  const embed = new EmbedBuilder()
    .setTitle(title)
    .setDescription(description)
    .setColor(config.color.error)
    .setTimestamp();
  await interaction.reply({
    embeds: [embed],
    flags: [MessageFlags.Ephemeral],
  });
};

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "general/tts" });
  if (!interaction.inGuild()) {
    await replyWithError(
      interaction,
      "Error - サーバー専用コマンド",
      "このコマンドはサーバー内でのみ使用できます。"
    );
    return;
  }

  if (!process.env.VOICEVOX_API_URL) {
    await replyWithError(
      interaction,
      "Error - VOICEVOX API URLが設定されていません",
      "このBotはTTSを使用できません。"
    );
    return;
  }

  const options = interaction.options as CommandInteractionOptionResolver;
  const group = options.getSubcommandGroup();
  const sub = options.getSubcommand();
  const command = group ? handlingCommands.get(group) : handlingCommands.get(sub);

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

    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) {
      errorChannel.send({ embeds: [logEmbed] });
    }

    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(interaction.client.config.color.error)
      .setTimestamp();

    if (interaction.channel) {
      await (interaction.channel as TextChannel).send({ embeds: [messageEmbed] });
    }

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
};
