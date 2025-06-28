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
  .setDescription("話者を選択");
export const execute = async (interaction: ChatInputCommandInteraction) => {
  await interaction.deferReply({
    flags: [MessageFlags.Ephemeral],
  });
  if (!RPC.rpc) {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
  }
  const speakers = await AudioLibrary.getSpeakers();
  if (speakers.length === 0) {
    const embed = new EmbedBuilder()
      .setTitle("Error: 利用可能な話者がいません")
      .setDescription("現在、利用可能な話者がいません。詳しくは管理者へお問い合わせください。")
      .setColor(interaction.client.config.color.error);
    return interaction.editReply({
      embeds: [embed],
    });
  }
  speakers.sort((a, b) => a.name.localeCompare(b.name));
  const components = [];
  if (speakers.length > 24) {
    const pagenation = [
      new ButtonBuilder()
        .setCustomId("tts_speaker_page_prev")
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
    components.push(new ActionRowBuilder<ButtonBuilder>().addComponents(pagenation));
    speakers.splice(24);
  }
  const selectMenu = new StringSelectMenuBuilder()
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
    ]);
  selectMenu.setMinValues(1).setMaxValues(1);
  components.push(new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(selectMenu));

  await interaction.editReply({
    content: "話者を選択してください。",
    components,
  });
};

export default {
  data,
  execute,
};
