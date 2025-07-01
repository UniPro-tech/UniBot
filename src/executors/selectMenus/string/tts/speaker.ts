import {
  ActionRowBuilder,
  EmbedBuilder,
  MessageFlags,
  StringSelectMenuBuilder,
  StringSelectMenuInteraction,
} from "discord.js";
import { uuid58Decode } from "@nakanoaas/uuid58";
import { AudioLibrary, RPC } from "voicevox.js";
import { writeTtsPreference } from "@/lib/dataUtils";

export const name = "speaker";

const connectVoicevox = async () => {
  if (!RPC.rpc) {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
  }
};

const handleSelect = async (interaction: StringSelectMenuInteraction) => {
  await connectVoicevox();
  const speakerId = uuid58Decode(interaction.values[0]);
  const speakerDetail = await AudioLibrary.getSpeaker(speakerId);

  const styleSelector = new StringSelectMenuBuilder()
    .setCustomId("tts_speaker_style")
    .setPlaceholder("スタイルを選択...")
    .addOptions(
      speakerDetail.styles.map((style) => ({
        label: style.name,
        value: style.id.toString(),
      }))
    );

  await interaction.editReply({
    content: `${speakerDetail.name} を選択しました。スタイルを選択してください。`,
    components: [new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(styleSelector)],
  });
};

const handleStyle = async (interaction: StringSelectMenuInteraction) => {
  const styleId = parseInt(interaction.values[0]);
  await connectVoicevox();

  const speakers = await AudioLibrary.getSpeakers();
  const speakerData = speakers.find((s) => s.styles.some((st) => st.id === styleId));
  const styleData = speakerData?.styles.find((st) => st.id === styleId);

  const embed = new EmbedBuilder()
    .setTitle("TTS話者設定")
    .setDescription("話者を設定しました。")
    .addFields(
      { name: "話者名", value: speakerData!.name, inline: true },
      { name: "スタイル名", value: styleData!.name, inline: true }
    )
    .setColor(interaction.client.config.color.success)
    .setFooter({
      text: "TTS設定",
      iconURL: interaction.client.user?.displayAvatarURL() || "",
    })
    .setTimestamp();

  await interaction.editReply({ embeds: [embed] });
  writeTtsPreference(interaction.user.id, "speaker", { styleId });
};

export const execute = async (interaction: StringSelectMenuInteraction) => {
  if (interaction.values[0] === "cancel") {
    await interaction.deferUpdate();
    await interaction.editReply({
      content: "設定をキャンセルしました。",
      components: [],
    });
    return;
  }

  await interaction.deferReply({ flags: [MessageFlags.Ephemeral] });

  const action = interaction.customId.split("_")[2];
  if (action === "select") {
    await handleSelect(interaction);
  } else if (action === "style") {
    await handleStyle(interaction);
  }
};
