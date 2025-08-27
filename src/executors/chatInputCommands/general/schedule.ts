import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";

export const handlingCommands = subCommandHandling("general/schedule");
export const data = addSubCommand(
  "general/schedule",
  new SlashCommandBuilder().setName("schedule").setDescription("予約投稿を管理します。")
);
export const guildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const subcommand = (interaction.options as CommandInteractionOptionResolver).getSubcommand();
  const command = handlingCommands.get(subcommand);

  if (!command) {
    console.info(`[Not Found] Command: ${subcommand}`);
    return;
  }

  try {
    await command.execute(interaction);
    console.info(`[Run] ${subcommand}`);
  } catch (error) {
    console.error(error);

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
