import { MessageFlags, StringSelectMenuInteraction } from "discord.js";
import { removeTtsDictionary } from "@/lib/dataUtils";

export const name = "dict";

export const execute = async (interaction: StringSelectMenuInteraction) => {
  const [, , action] = interaction.customId.split("_");
  const selectedValue = interaction.values[0];

  if (selectedValue === "cancel") {
    await interaction.reply({
      content: "選択がキャンセルされました。",
      flags: [MessageFlags.Ephemeral],
    });
    await interaction.message.delete().catch(() => {});
    return;
  }

  if (action === "remove") {
    console.log(`Removing TTS dictionary entry: ${selectedValue}`);
    removeTtsDictionary(undefined, undefined, undefined, selectedValue);
    await interaction.reply({
      content: "単語が削除されました。",
      flags: [MessageFlags.Ephemeral],
    });
  }
};
