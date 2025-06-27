import { uuid58Encode } from "@nakanoaas/uuid58";
import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChatInputCommandInteraction,
  SlashCommandSubcommandBuilder,
  StringSelectMenuBuilder,
} from "discord.js";
import { AudioLibrary, RPC } from "voicevox.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("speaker")
  .setDescription("Change the speaker of the TTS");
export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!RPC.rpc) {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
  }
  const speakers = await AudioLibrary.getSpeakers();
  if (speakers.length === 0) {
    return interaction.reply({
      content: "No speakers available. Please ask bot admin to add a speaker first.",
      ephemeral: true,
    });
  }
  speakers.sort((a, b) => a.name.localeCompare(b.name));
  const components = [];
  if (speakers.length > 24) {
    const pagenation = [
      new ButtonBuilder()
        .setCustomId("tts_speaker_page_prev_" + interaction.user.id)
        .setLabel("Previous")
        .setStyle(ButtonStyle.Primary)
        .setEmoji("⬅️")
        .setDisabled(true),
      new ButtonBuilder()
        .setCustomId("tts_speaker_page_next_" + interaction.user.id)
        .setLabel("Next")
        .setStyle(ButtonStyle.Primary)
        .setEmoji("➡️"),
    ];
    components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
    speakers.splice(24);
  }
  const selectMenu = new StringSelectMenuBuilder()
    .setCustomId("tts_speaker_select")
    .setPlaceholder("Select a speaker...")
    .addOptions([
      ...speakers.map((speaker) => ({
        label: speaker.name,
        value: uuid58Encode(speaker.speakerUuid),
      })),
      {
        label: "Cancel",
        value: "cancel",
        description: "Cancel the selection",
      },
    ]);
  selectMenu.setMinValues(1).setMaxValues(1);
  components.push(new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu));

  await interaction.reply({
    content: "Please select a speaker:",
    components,
  });
};

export default {
  data,
  execute,
};
