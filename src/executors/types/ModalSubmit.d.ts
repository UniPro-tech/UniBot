import { ModalSubmitInteraction } from "discord.js";

export interface ModalSubmitCommand {
  name: string;
  execute(interaction: ModalSubmitInteraction): Promise<void>;
}
