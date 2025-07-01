import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, MessageContextMenuCommandInteraction } from "discord.js";
import config from "@/config";

const MessageContextMenuCommandExecute = async (
  interaction: MessageContextMenuCommandInteraction
) => {
  const time = () => `[${interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}`;

  console.log(`${time()} info] MessageContextMenu ->${interaction.commandName}`);

  const command = interaction.client.interactionExecutorsCollections.messageContextMenuCommands.get(
    interaction.commandName
  );

  if (!command) {
    console.log(`${time()} info] Not Found: ${interaction.commandName}`);
    return;
  }

  if (!interaction.inGuild() && command.guildOnly) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはDMでは実行できません。")
      .setColor(interaction.client.config.color.error);

    await interaction.reply({ embeds: [embed] });
    console.log(`${time()} info] DM Only: ${interaction.commandName}`);
    return;
  }

  try {
    await command.execute(interaction);
    console.info(`${time()} run] ${interaction.commandName}`);

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
    console.error(
      `${time()} error]An Error Occurred in ${interaction.commandName}\nDetails:\n${error}`
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
