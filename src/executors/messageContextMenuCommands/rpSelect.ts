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
  .setName(name)
  .setType(ApplicationCommandType.Message);

export const execute = async (interaction: MessageContextMenuCommandInteraction) => {
  const targetMsg = interaction.targetMessage;
  const isRolePanel =
    targetMsg.author.id === interaction.client.user.id &&
    targetMsg.components[0]?.type === 1 &&
    targetMsg.components[0].components[0]?.customId?.startsWith("rp_");

  if (!isRolePanel) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("ロールパネルではありません")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [errorEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  await writeSelected({
    user: interaction.user.id,
    type: "Message",
    data: targetMsg.id,
  } as SelectedData);

  const successEmbed = new EmbedBuilder()
    .setTitle("ロールパネルを選択しました")
    .setColor(config.color.success)
    .setTimestamp();

  await interaction.reply({
    embeds: [successEmbed],
    flags: [MessageFlags.Ephemeral],
  });
};
