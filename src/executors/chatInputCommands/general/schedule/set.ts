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
  .setDescription("予約投稿を作成")
  .addBooleanOption((option) =>
    option.setName("repeat").setDescription("繰り返し投稿にするかどうか").setRequired(false)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (interaction.options.getBoolean("repeat")) {
    const modal = new ModalBuilder()
      .setCustomId("schedule_create_repeat")
      .setTitle("予約投稿の作成");
    const timeInput = new TextInputBuilder()
      .setCustomId("time")
      .setLabel("投稿時間")
      .setStyle(TextInputStyle.Short)
      .setPlaceholder(
        '例: "every day at 12:00"、 "every Monday at 09:00"\n詳しくはhelpコマンドを参照'
      )
      .setRequired(true);
    const contentInput = new TextInputBuilder()
      .setCustomId("message")
      .setLabel("投稿内容")
      .setStyle(TextInputStyle.Paragraph)
      .setPlaceholder("投稿内容を入力してください")
      .setRequired(true);

    const contentInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(contentInput);
    const timeInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(timeInput);
    const descriptionRow = new ActionRowBuilder<TextInputBuilder>().addComponents(
      new TextInputBuilder()
        .setCustomId("description")
        .setLabel("説明 (任意)")
        .setStyle(TextInputStyle.Paragraph)
        .setPlaceholder("この予約投稿の説明を入力してください (任意)")
        .setRequired(false)
    );

    modal.addComponents(timeInputRow, contentInputRow);
    await interaction.showModal(modal);
  } else {
    const modal = new ModalBuilder()
      .setCustomId("schedule_create_onetime")
      .setTitle("予約投稿の作成");
    const timeInput = new TextInputBuilder()
      .setCustomId("time")
      .setLabel("投稿時間 (YYYY-MM-DD HH:mm)")
      .setStyle(TextInputStyle.Short)
      .setPlaceholder("例: 2024-12-31 23:59")
      .setRequired(true);
    const contentInput = new TextInputBuilder()
      .setCustomId("message")
      .setLabel("投稿内容")
      .setStyle(TextInputStyle.Paragraph)
      .setPlaceholder("投稿内容を入力してください")
      .setRequired(true);

    const contentInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(contentInput);
    const timeInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(timeInput);

    modal.addComponents(timeInputRow, contentInputRow);
    await interaction.showModal(modal);
  }
};

export default {
  data,
  execute,
};
