import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
  TextChannel,
  MessageFlags,
} from "discord.js";
import { addSubCommand, addSubCommandGroup, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";

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
);

export const guildOnly = true;

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
    console.info(`[Not Found] Command: ${sub}`);
    return;
  }

  try {
    await command.execute(interaction);
    console.info(`[Run] ${sub}`);
  } catch (error) {
    console.error(error);

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
