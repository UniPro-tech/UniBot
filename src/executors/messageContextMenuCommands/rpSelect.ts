import { CommandInteraction, ContextMenuCommandBuilder, ApplicationCommandType } from "discord.js";

export const data = new ContextMenuCommandBuilder()
  .setName("RPを選択")
  .setType(3)
  .setType(ApplicationCommandType.Message);

export const adminGuildOnly = true;

export const execute = async (interaction: CommandInteraction) => {
  return;
};
