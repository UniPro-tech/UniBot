import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteractionOptionResolver,
  ChatInputCommandInteraction,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";

export const handlingCommands = subCommandHandling("admin/maintenance");
export const data = addSubCommand(
  "admin/maintenance",
  new SlashCommandBuilder().setName("maintenance").setDescription("メンテナンスモード")
);
export const guildOnly = true;
export const adminGuildOnly = true;

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const command = handlingCommands.get(
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

    await interaction.reply({ embeds: [messageEmbed] });
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
  return;
};
