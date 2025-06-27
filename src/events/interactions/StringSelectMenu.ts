import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, MessageFlags, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";

const StringSelectMenu = async (interaction: StringSelectMenuInteraction) => {
  try {
    const [prefix] = interaction.customId.split("_");
    const executionDefine =
      interaction.client.interactionExecutorsCollections.stringSelectMenus.get(prefix);
    if (!executionDefine) {
      console.log(
        `[${interaction.client.functions.timeUtils.timeToJSTstamp(
          Date.now(),
          true
        )} info] Not Found: ${interaction.customId}`
      );
      return;
    }
    console.log(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} info] StringSelectMenu -> ${interaction.customId}`
    );
    await executionDefine.execute(interaction);
  } catch (error) {
    console.error(
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured in ${interaction.customId}\nDetails:\n${error}`
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
    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({ embeds: [messageEmbed] });
    }
    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) {
      logChannel.send({ embeds: [messageEmbed] });
    }
  }
};

export default StringSelectMenu;
