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

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.inGuild()) {
    const embed = new EmbedBuilder()
      .setTitle("Error - サーバー専用コマンド")
      .setDescription("このコマンドはサーバー内でのみ使用できます。")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [embed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  if (!process.env.VOICEVOX_API_URL) {
    const embed = new EmbedBuilder()
      .setTitle("Error - VOICEVOX API URLが設定されていません")
      .setDescription("このBotはTTSを使用できません。")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [embed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  const command =
    (interaction.options as CommandInteractionOptionResolver).getSubcommandGroup() != null
      ? handlingCommands.get(
          (interaction.options as CommandInteractionOptionResolver).getSubcommandGroup() as string
        )
      : handlingCommands.get(
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
    console.error(error);
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

    await (interaction.channel as TextChannel).send({ embeds: [messageEmbed] });
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
  return;
};
