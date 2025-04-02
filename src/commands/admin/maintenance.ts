import {
  SlashCommandBuilder,
  EmbedBuilder,
  CommandInteraction,
  CommandInteractionOptionResolver,
} from "discord.js";
import { addSubCommand, subCommandHandling } from "@/lib/commandUtils";
import { GetLogChannel, GetErrorChannel } from "@/lib/channelUtils";
import config from "@/config";

export const handlingCommands = subCommandHandling("admin/maintenance");
export const data = addSubCommand(
  "admin/maintenance",
  new SlashCommandBuilder()
    .setName("maintenance")
    .setDescription("メンテナンスモード")
);
export const guildOnly = true;

export const execute = async (interaction: CommandInteraction) => {
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
      `[Run] ${(
        interaction.options as CommandInteractionOptionResolver
      ).getSubcommand()}`
    );

    const logEmbed = new EmbedBuilder()
      .setTitle("サブコマンド実行ログ")
      .setDescription(`${interaction.user} がサブコマンドを実行しました。`)
      .setColor(interaction.client.config.color.s)
      .setTimestamp()
      .setThumbnail(interaction.user.displayAvatarURL())
      .addFields([
        {
          name: "サブコマンド",
          value: `\`\`\`\n/${(
            interaction.options as CommandInteractionOptionResolver
          ).getSubcommand()}\n\`\`\``,
        },
        {
          name: "実行サーバー",
          value:
            "```\n" + interaction.inGuild()
              ? `${interaction.guild?.name} (${interaction.guild?.id})`
              : "DM" + "\n```",
        },
        {
          name: "実行ユーザー",
          value:
            "```\n" +
            `${interaction.user.tag}(${interaction.user.id})` +
            "\n```",
        },
      ])
      .setFooter({ text: `${interaction.id}` });
    const channel = await GetLogChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
  } catch (error) {
    console.error(error);
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + (error as any).toString() + "\n```")
      .setColor(config.color.e)
      .setTimestamp();

    const channel = await GetErrorChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません、エラーが発生しました...")
      .setDescription("```\n" + error + "\n```")
      .setColor(interaction.client.config.color.e)
      .setTimestamp();

    await interaction.reply({ embeds: [messageEmbed] });
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
  return;
};
