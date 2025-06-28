import {
  ActionRowBuilder,
  MessageFlags,
  StringSelectMenuBuilder,
  StringSelectMenuInteraction,
} from "discord.js";
import { uuid58Decode } from "@nakanoaas/uuid58";
import { AudioLibrary, RPC } from "voicevox.js";
import { writeTtsPreference } from "@/lib/dataUtils";

export const name = "speaker";

export const execute = async (interaction: StringSelectMenuInteraction) => {
  if (interaction.values[0] === "cancel") {
    await interaction.deferUpdate();
    await interaction.editReply({
      content: "設定をキャンセルしました。",
      components: [],
    });
    return;
  }
  await interaction.deferReply({
    flags: [MessageFlags.Ephemeral],
  });
  switch (interaction.customId.split("_")[2]) {
    case "select":
      if (!RPC.rpc) {
        const headers = {
          Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
        };
        await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
      }
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
      break;
    case "style":
      const styleId = parseInt(interaction.values[0]);
      if (!RPC.rpc) {
        const headers = {
          Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
        };
        await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
      }
      const speakers = await AudioLibrary.getSpeakers();
      const speakerData = speakers.find((s) => s.styles.find((st) => st.id === styleId));
      await interaction.editReply({
        content: `${speakerData!.name} の ${
          speakerData!.styles.find((st) => st.id === styleId)!.name
        } を設定しました。`,
      });
      writeTtsPreference(interaction.user.id, "speaker", {
        styleId,
      });
      break;
  }
};
