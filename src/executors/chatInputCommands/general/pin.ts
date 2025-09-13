import { SlashCommandBuilder } from "@discordjs/builders";
import {
  ActionRowBuilder,
  ChatInputCommandInteraction,
  ModalBuilder,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";

export const guildOnly = false;

export const data = new SlashCommandBuilder()
  .setName("pin")
  .setDescription("メッセージをピン留めします。");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.guild) {
    await interaction.reply({
      content: "このコマンドはサーバー内でのみ使用できます。",
      ephemeral: true,
    });
    return "No data";
  }
  if (!interaction.memberPermissions?.has("ManageMessages")) {
    await interaction.reply({
      content: "このコマンドを使用する権限がありません。",
      ephemeral: true,
    });
    return "No data";
  }
  const modal = new ModalBuilder().setCustomId("pin_message").setTitle("メッセージのピン留め");
  const contentInput = new TextInputBuilder()
    .setCustomId("message")
    .setLabel("投稿内容")
    .setStyle(TextInputStyle.Paragraph)
    .setPlaceholder("投稿内容を入力してください")
    .setRequired(true);

  const contentInputRow = new ActionRowBuilder<TextInputBuilder>().addComponents(contentInput);

  modal.addComponents(contentInputRow);
  await interaction.showModal(modal);
  return "No data";
};

export default {
  guildOnly,
  data,
  execute,
};
