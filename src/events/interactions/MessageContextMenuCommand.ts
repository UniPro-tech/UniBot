import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, MessageContextMenuCommandInteraction } from "discord.js";
import config from "@/config";

const MessageContextMenuCommandExecute = async (
  interaction: MessageContextMenuCommandInteraction
) => {
  console.log(
    `[${interaction.client.functions.timeUtils.timeToJSTstamp(
      Date.now(),
      true
    )} info] MessageContextMenu ->${interaction.commandName}`
  );
  const command = interaction.client.interactionExecutorsCollections.messageContextMenuCommands.get(
    interaction.commandName
  );
  if (!command) {
    console.log(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} info] Not Found: ${interaction.commandName}`
    );
    return;
  }
  if (!interaction.inGuild() && command.guildOnly) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはDMでは実行できません。")
      .setColor(interaction.client.config.color.error);
    interaction.reply({ embeds: [embed] });
    console.log(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true)} info] DM Only: ${
        interaction.commandName
      }`
    );
    return;
  }

  try {
    await command.execute(interaction);
    console.info(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true)} run] ${
        interaction.commandName
      }`
    );

    const logEmbed = new EmbedBuilder()
      .setTitle("コマンド実行ログ")
      .setDescription(`${interaction.user} がコマンドを実行しました。`)
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
    const channel = await GetLogChannel(interaction.client);
    if (channel) {
      channel.send({ embeds: [logEmbed] });
    }
  } catch (error) {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred in ${interaction.commandName}\nDetails:\n${error}`
    );
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
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    await interaction.reply({ embeds: [messageEmbed] });
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
};

export default MessageContextMenuCommandExecute;
