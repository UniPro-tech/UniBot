import { uuid58Encode } from "@nakanoaas/uuid58";
import {
  ActionRowBuilder,
  ButtonBuilder,
  ButtonStyle,
  ChatInputCommandInteraction,
  EmbedBuilder,
  MessageFlags,
  SlashCommandSubcommandBuilder,
  StringSelectMenuBuilder,
} from "discord.js";
import { AudioLibrary, RPC } from "voicevox.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("speaker")
  .setDescription("TTSの話者を選択");

const createErrorEmbed = (color: number) =>
  new EmbedBuilder()
    .setTitle("Error: 利用可能な話者がいません")
    .setDescription("現在、利用可能な話者がいません。詳しくは管理者へお問い合わせください。")
    .setColor(color);

const createPaginationButtons = () => [
  new ButtonBuilder()
    .setCustomId("tts_speaker_page_prev_1")
    .setLabel("Previous")
    .setStyle(ButtonStyle.Primary)
    .setEmoji("⬅️")
    .setDisabled(true),
  new ButtonBuilder()
    .setCustomId("tts_speaker_page_next_2")
    .setLabel("Next")
    .setStyle(ButtonStyle.Primary)
    .setEmoji("➡️"),
];

const createSpeakerSelectMenu = (speakers: any[]) =>
  new StringSelectMenuBuilder()
    .setCustomId("tts_speaker_select")
    .setPlaceholder("話者を選択...")
    .addOptions([
      ...speakers.map((speaker) => ({
        label: speaker.name,
        value: uuid58Encode(speaker.speakerUuid),
      })),
      {
        label: "キャンセル",
        value: "cancel",
        description: "選択をキャンセル",
      },
    ])
    .setMinValues(1)
    .setMaxValues(1);

export const execute = async (interaction: ChatInputCommandInteraction) => {
  await interaction.deferReply({ flags: [MessageFlags.Ephemeral] });

  if (!RPC.rpc) {
    await RPC.connect(process.env.VOICEVOX_API_URL as string, {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    });
  }

  const speakers = await AudioLibrary.getSpeakers();
  if (!speakers.length) {
    const embed = createErrorEmbed(interaction.client.config.color.error);
    return interaction.editReply({ embeds: [embed] });
  }

  speakers.sort((a, b) => a.name.localeCompare(b.name));
  const components = [];

  if (speakers.length > 24) {
    components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(createPaginationButtons()));
    speakers.splice(24);
  }

  components.push(
    new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(createSpeakerSelectMenu(speakers))
  );

  await interaction.editReply({
    content: "話者を選択してください。",
    components,
  });
};

export default { data, execute };
