import {
  ChatInputCommandInteraction,
  CommandInteractionOptionResolver,
  EmbedBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";

export const data = addSubCommand(
  "general/tts/group/dict",
  new SlashCommandSubcommandGroupBuilder().setName("dict").setDescription("TTSの辞書を管理")
);

export const handlingCommand = subCommandHandling("general/tts/group/dict");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const subcommand = (interaction.options as CommandInteractionOptionResolver).getSubcommand();
  const command = handlingCommand.get(subcommand);

  if (!command) {
    console.info(`[Not Found] Command: ${subcommand}`);
    return;
  }

  try {
    await command.execute(interaction);
    console.info(`[Run] ${subcommand}`);
  } catch (error) {
    console.error(`[Error] ${subcommand}`, error);

    const embed = new EmbedBuilder()
      .setTitle("TTS Command Error")
      .setDescription(`Error executing command: ${error}`)
      .setColor(interaction.client.config.color.error);

    const [logChannel, errorChannel] = await Promise.all([
      GetLogChannel(interaction.client),
      GetErrorChannel(interaction.client),
    ]);

    if (logChannel) await logChannel.send({ embeds: [embed] });
    if (errorChannel) await errorChannel.send({ embeds: [embed] });
  }
};
