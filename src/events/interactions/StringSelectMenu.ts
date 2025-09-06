import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { EmbedBuilder, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";
import { loggingSystem } from "@/index";

const StringSelectMenu = async (interaction: StringSelectMenuInteraction) => {
  const logger = loggingSystem.getLogger({ function: "StringSelectMenu" });

  try {
    const [prefix] = interaction.customId.split("_");
    const executor =
      interaction.client.interactionExecutorsCollections.stringSelectMenus.get(prefix);

    if (!executor) {
      logger.error(
        { extra_context: { customId: interaction.customId } },
        "No executor found for this StringSelectMenu"
      );
      return;
    }

    logger.info(
      { extra_context: { customId: interaction.customId } },
      "StringSelectMenu execution started"
    );
    await executor.execute(interaction);
  } catch (error) {
    const errorMsg = (error as Error).toString();
    logger.error(
      { extra_context: { customId: interaction.customId }, stack_trace: (error as Error).stack },
      "Error executing StringSelectMenu",
      error
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
