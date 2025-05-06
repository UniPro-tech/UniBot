import { StringSelectMenuInteraction } from "discord.js";

export interface StringSelectMenuDefineType {
  name: string;
  execute(interaction: StringSelectMenuInteraction, id: string): Promise<void>;
}
