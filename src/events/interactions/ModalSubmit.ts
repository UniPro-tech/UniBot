import { ALStorage, loggingSystem } from "@/index";
import { EmbedBuilder, ModalSubmitInteraction } from "discord.js";

const ModalSubmitExecute = async (interaction: ModalSubmitInteraction) => {
  const ctx = {
    ...ALStorage.getStore(),
    user_id: interaction.user.id,
    context: { discord: { guild: interaction.guild?.id, channel: interaction.channel?.id } },
  };
  const logger = loggingSystem.getLogger({ ...ctx, function: "ModalSubmitExecute" });
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
    ALStorage.run(ctx, async () => {
      await modal.execute(interaction);
    });
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
