import {
  ContextMenuCommandBuilder,
  ApplicationCommandType,
  MessageContextMenuCommandInteraction,
} from "discord.js";

export const name = "RPを選択";

export const data = new ContextMenuCommandBuilder()
  .setName("RPを選択")
  .setType(3)
  .setType(ApplicationCommandType.Message);

export const execute = async (interaction: MessageContextMenuCommandInteraction) => {
  interaction.reply("OK");
  return;
};
