import { EmbedBuilder, ModalSubmitInteraction } from "discord.js";

const ModalSubmitExecute = async (interaction: ModalSubmitInteraction) => {
  try {
    const time = () =>
      `[${interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}`;

    console.log(`${time()} info] ModalSubmit ->${interaction.customId}`);

    const modal = interaction.client.interactionExecutorsCollections.modalSubmitCommands.get(
      interaction.customId
    );
    if (!modal) {
      console.log(
        `${time()} error] ModalSubmit -> No modal found for customId: ${interaction.customId}`
      );
      return;
    }
    modal.execute(interaction);
  } catch (error) {
    console.error("Error handling modal submit interaction:", error);
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
