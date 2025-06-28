import { Events, Interaction } from "discord.js";
import ChatInputCommandExecute from "@/events/interactions/ChatInputCommand";
import StringSelectMenuExecute from "@/events/interactions/StringSelectMenu";
import MessageContextMenuCommandExecute from "./interactions/MessageContextMenuCommand";
import ButtonExecute from "./interactions/Button";

export const name = Events.InteractionCreate;
export const execute = async (interaction: Interaction) => {
  if (interaction.isChatInputCommand()) {
    await ChatInputCommandExecute(interaction);
  } else if (interaction.isStringSelectMenu()) {
    await StringSelectMenuExecute(interaction);
  } else if (interaction.isMessageContextMenuCommand()) {
    await MessageContextMenuCommandExecute(interaction);
  } else if (interaction.isButton()) {
    await ButtonExecute(interaction);
  }
};

export default {
  name,
  execute,
};
