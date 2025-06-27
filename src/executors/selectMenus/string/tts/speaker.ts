import {
  ActionRowBuilder,
  MessageFlags,
  StringSelectMenuBuilder,
  StringSelectMenuInteraction,
} from "discord.js";
import { uuid58Decode } from "@nakanoaas/uuid58";
import { AudioLibrary, RPC } from "voicevox.js";
import { write } from "fs";
import { writeTtsPreference } from "@/lib/dataUtils";

export const name = "speaker";

export const execute = async (interaction: StringSelectMenuInteraction) => {
  if (interaction.values[0] === "cancel") {
    await interaction.reply({
      content: "Selection cancelled.",
      flags: [MessageFlags.Ephemeral],
    });
    await interaction.message.delete().catch(() => {});
    return;
  }
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
        .setPlaceholder("Select a style...")
        .addOptions(
          speakerDetail.styles.map((style) => ({
            label: style.name,
            value: style.id.toString(),
          }))
        );
      await interaction.reply({
        components: [new ActionRowBuilder<StringSelectMenuBuilder>().addComponents(styleSelector)],
      });
      break;
    case "style":
      const styleId = parseInt(interaction.values[0]);
      await interaction.reply({
        content: `You selected style ${styleId}.`,
        flags: [MessageFlags.Ephemeral],
      });
      writeTtsPreference(interaction.user.id, "speaker", {
        styleId,
      });
      break;
  }
};
