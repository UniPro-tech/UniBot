import { writeTtsDictionary } from "@/lib/dataUtils";
import { PrismaClientKnownRequestError } from "@prisma/client/runtime/library";
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

const createEmbed = (
  title: string,
  description: string,
  color: number,
  fields?: { name: string; value: string; inline?: boolean }[]
) => {
  const embed = new EmbedBuilder().setTitle(title).setColor(color);
  if (description) embed.setDescription(description);
  if (fields) embed.addFields(fields);
  return embed;
};

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.guild) {
    const embed = createEmbed(
      "エラー",
      "このコマンドはサーバー内でのみ使用できます。",
      interaction.client.config.color.error
    );
    await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    return;
  }

  const word = interaction.options.getString("word", true);
  const definition = interaction.options.getString("definition", true);
  const caseSensitive = interaction.options.getBoolean("case_sensitive") ?? false;

  try {
    await writeTtsDictionary(
      interaction.user.id,
      interaction.guild.id,
      word,
      definition,
      caseSensitive
    );

    const embed = createEmbed(
      "単語を辞書に追加しました！",
      "",
      interaction.client.config.color.success,
      [
        { name: "単語", value: word, inline: true },
        { name: "読み", value: definition, inline: true },
      ]
    );
    await interaction.reply({ embeds: [embed] });
  } catch (error) {
    if ((error as PrismaClientKnownRequestError).code === "P2002") {
      const embed = createEmbed(
        "エラー",
        "この単語はすでに辞書に存在します。",
        interaction.client.config.color.error
      );
      await interaction.reply({ embeds: [embed], flags: [MessageFlags.Ephemeral] });
    } else {
      throw error;
    }
  }
};

export default {
  data,
  execute,
};
