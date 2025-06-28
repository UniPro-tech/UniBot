import { listTtsDictionary } from "@/lib/dataUtils";
import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChatInputCommandInteraction,
  EmbedBuilder,
  MessageFlags,
  PermissionsBitField,
  SlashCommandSubcommandBuilder,
  StringSelectMenuBuilder,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("remove")
  .setDescription("Remove a word from the dictionary");
export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (interaction.guild === null) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはサーバー内でのみ使用できます。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }
  await interaction.deferReply();
  const allWords = await listTtsDictionary(
    interaction.guild.id,
    !(interaction.member?.permissions as PermissionsBitField).has(
      PermissionsBitField.Flags.Administrator
    )
      ? interaction.user.id
      : undefined
  );
  if (allWords.length === 0) {
    return interaction.editReply({
      content: "辞書に登録されている単語がありません。",
    });
  }
  const components = [];
  if (allWords.length > 24) {
    const buttons = [
      new ButtonBuilder()
        .setCustomId("tts_dict_remove_page_prev_1")
        .setDisabled(true)
        .setLabel("Previous")
        .setEmoji("◀️")
        .setStyle(ButtonStyle.Primary),
      new ButtonBuilder()
        .setCustomId("tts_dict_remove_page_next_1")
        .setLabel("Next")
        .setEmoji("▶️")
        .setStyle(ButtonStyle.Primary),
    ];
    const row = new ActionRowBuilder<ButtonBuilder>().addComponents(buttons);
    components.push(row);
    allWords.splice(25);
  }
  const selectMenu = new StringSelectMenuBuilder()
    .setCustomId("tts_dict_remove")
    .setPlaceholder("削除する単語を選択...")
    .addOptions(
      allWords.map((word) => ({
        label: word.word,
        value: word.id,
        description: word.definition,
      }))
    );
  const selectRow = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu);
  components.push(selectRow);

  await interaction.editReply({
    content: "削除する単語を選択してください。",
    components: components,
  });
};

export default {
  data,
  execute,
};
