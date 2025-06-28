import { writeTtsDictionary } from "@/lib/dataUtils";
import { ChatInputCommandInteraction, SlashCommandSubcommandBuilder } from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("add")
  .setDescription("Add a new word to the dictionary")
  .addStringOption((option) =>
    option.setName("word").setDescription("The word to add").setRequired(true)
  )
  .addStringOption((option) =>
    option.setName("definition").setDescription("The definition of the word").setRequired(true)
  );
export const execute = async (interaction: ChatInputCommandInteraction) => {
  const word = interaction.options.getString("word");
  const definition = interaction.options.getString("definition");

  if (!word || !definition) {
    await interaction.reply("Please provide both a word and its definition.");
    return;
  }

  if (interaction.guild === null) {
    await interaction.reply("This command can only be used in a server.");
    return;
  }

  await writeTtsDictionary(interaction.user.id, interaction.guild.id, word, definition);

  await interaction.reply(`Added "${word}" to the dictionary with definitions: ${definition}`);
};

export default {
  data,
  execute,
};
