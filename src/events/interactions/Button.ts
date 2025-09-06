import { GetErrorChannel, GetLogChannel } from "@/lib/channelUtils";
import { ButtonInteraction, EmbedBuilder } from "discord.js";
import config from "@/config";
import { loggingSystem } from "@/index";

const ButtonExecute = async (interaction: ButtonInteraction) => {
  const logger = loggingSystem.getLogger({ function: "ButtonExecute" });
  try {
    const [prefix] = interaction.customId.split("_");
    const executionDefine = interaction.client.interactionExecutorsCollections.buttons.get(prefix);

    if (!executionDefine) {
      logger.error(
        { extra_context: { customId: interaction.customId } },
        "No execution definition found"
      );
      return;
    }

    logger.info(
      { extra_context: { customId: interaction.customId } },
      "Button interaction executed"
    );
    await executionDefine.execute(interaction);
  } catch (error) {
    const errorMsg = (error as Error).toString();
    logger.error(
      { extra_context: { customId: interaction.customId }, stack_trace: (error as Error).stack },
      "Button interaction execution failed",
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

    const [errorChannel, logChannel] = await Promise.all([
      GetErrorChannel(interaction.client),
      GetLogChannel(interaction.client),
    ]);

    if (errorChannel) errorChannel.send({ embeds: [logEmbed] });
    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({ embeds: [messageEmbed] });
    }
    if (logChannel) logChannel.send({ embeds: [messageEmbed] });
  }
};

export default ButtonExecute;
