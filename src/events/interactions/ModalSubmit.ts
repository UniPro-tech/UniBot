import { loggingSystem } from "@/index";
import { EmbedBuilder, ModalSubmitInteraction } from "discord.js";

const ModalSubmitExecute = async (interaction: ModalSubmitInteraction) => {
  const logger = loggingSystem.getLogger({ function: "ModalSubmitExecute" });
  try {
    logger.info(
      { extra_context: { customId: interaction.customId } },
      "ModalSubmit execution started"
    );

    const modal = interaction.client.interactionExecutorsCollections.modalSubmitCommands.get(
      interaction.customId
    );
    if (!modal) {
      logger.error(
        { service: "ModalSubmitExecutor", customId: interaction.customId },
        "Modal not found"
      );
      return;
    }
    modal.execute(interaction);
  } catch (error) {
    logger.error(
      { extra_context: { customId: interaction.customId }, stack_trace: (error as Error).stack },
      "Error executing ModalSubmit",
      error
    );
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("モーダルの送信処理中にエラーが発生しました。")
      .setColor("Red")
      .setTimestamp();

    if (interaction.replied || interaction.deferred) {
      await interaction.followUp({ embeds: [embed], ephemeral: true });
    } else {
      await interaction.reply({ embeds: [embed], ephemeral: true });
    }
  }
};
export default ModalSubmitExecute;
