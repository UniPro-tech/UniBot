import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";

const StringSelectMenu = async (interaction: StringSelectMenuInteraction) => {
  const now = Date.now();
  const time = interaction.client.functions.timeUtils.timeToJSTstamp(now, true);

  try {
    const [prefix] = interaction.customId.split("_");
    const executor =
      interaction.client.interactionExecutorsCollections.stringSelectMenus.get(prefix);

    if (!executor) {
      console.log(`[${time} info] Not Found: ${interaction.customId}`);
      return;
    }

    console.log(`[${time} info] StringSelectMenu -> ${interaction.customId}`);
    await executor.execute(interaction);
  } catch (error) {
    const errorMsg = (error as Error).toString();
    console.error(
      `[${time} error] An Error Occured in ${interaction.customId}\nDetails:\n${errorMsg}`
    );

    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const messageEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription(`\`\`\`\n${errorMsg}\n\`\`\``)
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) errorChannel.send({ embeds: [logEmbed] });

    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({ embeds: [messageEmbed] });
    }

    const logChannel = await GetLogChannel(interaction.client);
    if (logChannel) logChannel.send({ embeds: [messageEmbed] });
  }
};

export default StringSelectMenu;
