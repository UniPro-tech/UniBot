import { SlashCommandBuilder } from "@discordjs/builders";
import {
  ActionRowBuilder,
  ChatInputCommandInteraction,
  InteractionContextType,
  ModalBuilder,
  PermissionFlagsBits,
  TextInputBuilder,
  TextInputStyle,
} from "discord.js";

export const guildOnly = false;

export const data = new SlashCommandBuilder()
  .setName("pin")
  .setDescription("メッセージをピン留めします。")
  .setDefaultMemberPermissions(PermissionFlagsBits.PinMessages)
  .setContexts(InteractionContextType.Guild);

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const modal = new ModalBuilder().setCustomId("pin_message").setTitle("メッセージのピン留め");
  const contentInput = new TextInputBuilder()
    .setCustomId("message")
    .setLabel("投稿内容")
    .setStyle(TextInputStyle.Paragraph)
    .setPlaceholder(
      "投稿内容を入力してください。すでにPinされたメッセージがある場合は上書きされます。"
    )
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
