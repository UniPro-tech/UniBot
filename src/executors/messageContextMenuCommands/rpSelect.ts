import {
  ContextMenuCommandBuilder,
  ApplicationCommandType,
  MessageContextMenuCommandInteraction,
  MessageFlags,
  EmbedBuilder,
} from "discord.js";
import { SelectedData, writeSelected } from "@/lib/dataUtils.js";
import config from "@/config.js";

export const name = "RPを選択";

export const data = new ContextMenuCommandBuilder()
  .setName("RPを選択")
  .setType(3)
  .setType(ApplicationCommandType.Message);

export const execute = async (interaction: MessageContextMenuCommandInteraction) => {
  if (
    interaction.targetMessage.author.id != interaction.client.user.id ||
    !(
      interaction.targetMessage.components[0] &&
      interaction.targetMessage.components[0].type == 1 &&
      interaction.targetMessage.components[0].components[0].customId?.startsWith("rp_")
    )
  ) {
    const messageEmbed = new EmbedBuilder()
      .setTitle("ロールパネルではありません")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [messageEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  const messageId = interaction.targetMessage.id;
  writeSelected({
    user: interaction.user.id,
    type: "Message",
    data: interaction.targetMessage.id,
  } as SelectedData);
  const messageEmbed = new EmbedBuilder()
    .setTitle("ロールパネルを選択しました")
    .setColor(config.color.success)
    .setTimestamp();
  await interaction.reply({
    embeds: [messageEmbed],
    flags: [MessageFlags.Ephemeral],
  });
  return;
};
