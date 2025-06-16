import { StringSelectMenuInteraction } from "discord.js";

export interface StringSelectMenu {
  name: string;
  execute(interaction: StringSelectMenuInteraction, id: string): Promise<void>;
}
