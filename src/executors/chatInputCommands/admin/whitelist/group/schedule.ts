import {
  ChatInputCommandInteraction,
  CommandInteractionOptionResolver,
  EmbedBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";

export const data = addSubCommand(
  "admin/whitelist/group/schedule",
  new SlashCommandSubcommandGroupBuilder()
    .setName("schedule")
    .setDescription("スケジュールコマンドのホワイトリスト設定を変更")
);

export const handlingCommand = subCommandHandling("admin/whitelist/group/schedule");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const options = interaction.options as CommandInteractionOptionResolver;
  const subcommand = options.getSubcommand();
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
      .setTitle("Whitelist Command Error")
      .setDescription(`Error executing command: ${error}`)
      .setColor("Red");

    const [logChannel, errorChannel] = await Promise.all([
      GetLogChannel(interaction.client),
      GetErrorChannel(interaction.client),
    ]);

    if (logChannel) await logChannel.send({ embeds: [embed] });
    if (errorChannel) await errorChannel.send({ embeds: [embed] });
  }
};
