import { MessageFlags, StringSelectMenuInteraction } from "discord.js";
import { removeTtsDictionary } from "@/lib/dataUtils";

export const name = "dict";

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
    case "remove":
      removeTtsDictionary(undefined, undefined, undefined, interaction.values[0]);
      await interaction.reply({
        content: "Word removed from the dictionary.",
        flags: [MessageFlags.Ephemeral],
      });
  }
};
