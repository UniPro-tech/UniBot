import { listTtsDictionary } from "@/lib/dataUtils";
import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChatInputCommandInteraction,
  MessageFlags,
  PermissionsBitField,
  SlashCommandSubcommandBuilder,
  StringSelectMenuBuilder,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("remove")
  .setDescription("Remove a word from the dictionary");
export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.guild) {
    return interaction.reply("This command can only be used in a server.");
  }
  const allWords = await listTtsDictionary(
    interaction.guild.id,
    !(interaction.member?.permissions as PermissionsBitField).has(
      PermissionsBitField.Flags.Administrator
    )
      ? interaction.user.id
      : undefined
  );
  if (allWords.length === 0) {
    return interaction.reply("No words found in the dictionary.");
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
    .setPlaceholder("Select a word to remove")
    .addOptions(
      allWords.map((word) => ({
        label: word.word,
        value: word.id,
      }))
    );
  const selectRow = new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu);
  components.push(selectRow);

  await interaction.reply({
    content: "Please select a word to remove:",
    components: components,
    flags: [MessageFlags.Ephemeral],
  });
};

export default {
  data,
  execute,
};
