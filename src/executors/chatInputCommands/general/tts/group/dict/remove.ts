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
  .setDescription("TTS辞書から単語を削除");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.guild) {
    const embed = new EmbedBuilder()
      .setTitle("Error - サーバー専用コマンド")
      .setDescription("このコマンドはサーバー内でのみ使用できます。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }

  await interaction.deferReply({ flags: [MessageFlags.Ephemeral] });

  const isAdmin = (interaction.member?.permissions as PermissionsBitField)?.has(
    PermissionsBitField.Flags.Administrator
  );
  const allWords = await listTtsDictionary(
    interaction.guild.id,
    isAdmin ? undefined : interaction.user.id
  );

  if (!allWords.length) {
    await interaction.editReply({
      content: "辞書に登録されてる単語がありません。",
    });
    return;
  }

  const components = [];

  if (allWords.length > 24) {
    const prevBtn = new ButtonBuilder()
      .setCustomId("tts_dict_remove_page_prev_1")
      .setDisabled(true)
      .setLabel("Previous")
      .setEmoji("◀️")
      .setStyle(ButtonStyle.Primary);

    const nextBtn = new ButtonBuilder()
      .setCustomId("tts_dict_remove_page_next_1")
      .setLabel("Next")
      .setEmoji("▶️")
      .setStyle(ButtonStyle.Primary);

    components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(prevBtn, nextBtn));
    allWords.splice(25);
  }

  const selectMenu = new StringSelectMenuBuilder()
    .setCustomId("tts_dict_remove")
    .setPlaceholder("削除する単語を選んでください。")
    .addOptions(
      allWords.map((word) => ({
        label: word.word,
        value: word.id,
        description: word.definition,
      }))
    );

  components.push(new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu));

  await interaction.editReply({
    content: "削除したい単語を選んでください。",
    components,
  });
};

export default {
  data,
  execute,
};
