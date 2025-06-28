import { StringSelectMenuInteraction } from "discord.js";

export interface StringSelectMenu {
  name: string;
  execute(interaction: StringSelectMenuInteraction): Promise<void>;
}
