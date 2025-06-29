import { writeTtsDictionary } from "@/lib/dataUtils";
import {
  ChatInputCommandInteraction,
  EmbedBuilder,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("add")
  .setDescription("TTS辞書に単語を追加")
  .addStringOption((option) =>
    option.setName("word").setDescription("追加する単語").setRequired(true)
  )
  .addStringOption((option) =>
    option.setName("definition").setDescription("追加する単語の読み").setRequired(true)
  )
  .addBooleanOption((option) =>
    option
      .setName("case_sensitive")
      .setDescription("大文字小文字の区別をするか (デフォルト: false)")
      .setRequired(false)
  );
export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (interaction.guild === null) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはサーバー内でのみ使用できます。")
      .setColor(interaction.client.config.color.error);
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }

  const word = interaction.options.getString("word");
  const definition = interaction.options.getString("definition");

  await writeTtsDictionary(
    interaction.user.id,
    interaction.guild.id,
    word!,
    definition!,
    interaction.options.getBoolean("case_sensitive") ?? false
  );

  const embed = new EmbedBuilder()
    .setTitle("単語を辞書に追加しました！")
    .addFields([
      { name: "単語", value: word!, inline: true },
      { name: "読み", value: definition!, inline: true },
    ])
    .setColor(interaction.client.config.color.success);
  await interaction.reply({ embeds: [embed] });
};

export default {
  data,
  execute,
};
