import {
  ActionRowBuilder,
  ChatInputCommandInteraction,
  ModalBuilder,
  SlashCommandSubcommandBuilder,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("set")
  .setDescription("予約投稿を作成");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const modal = new ModalBuilder().setCustomId("schedule_create").setTitle("予約投稿の作成");
  const timeInput = new TextInputBuilder()
    .setCustomId("time")
    .setLabel("投稿時間 (YYYY-MM-DD HH:mm)")
    .setStyle(TextInputStyle.Short)
    .setPlaceholder("例: 2024-12-31 23:59")
    .setRequired(true);
  const contentInput = new TextInputBuilder()
    .setCustomId("content")
    .setLabel("投稿内容")
    .setStyle(TextInputStyle.Paragraph)
    .setPlaceholder("投稿内容を入力してください")
    .setRequired(true);

  const contentInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(contentInput);
  const timeInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(timeInput);

  modal.addComponents(timeInputRow, contentInputRow);
  await interaction.showModal(modal);
};

export default {
  data,
  execute,
};
