import config from "@/config";
import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { writeChatInputCommandInteractionLog } from "@/lib/logger";
import { ChatInputCommandInteraction, EmbedBuilder } from "discord.js";

const ChatInputCommandExecute = async (interaction: ChatInputCommandInteraction) => {
  console.log(
    `[${interaction.client.functions.timeUtils.timeToJSTstamp(
      Date.now(),
      true
    )} info] ChatInputCommand->${interaction.commandName}`
  );
  const command = interaction.client.interactionExecutorsCollections.chatInputCommands.get(
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
    console.log(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true)} run] ${
        interaction.commandName
      }`
    );
    await writeChatInputCommandInteractionLog(interaction);
  } catch (error) {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured in ${interaction.commandName}\nDatails:\n${error}`
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

export default ChatInputCommandExecute;
