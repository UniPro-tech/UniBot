import { ChatInputCommandInteraction, SlashCommandBuilder } from "discord.js";

export interface Command {
  handlingCommands: Function | undefined;
  admunGuildOnly: boolean | undefined;
  guildOnly: boolean | undefined;
  data:
    | SlashCommandBuilder
    | Omit<SlashCommandBuilder, "addSubcommandGroup" | "addSubcommand">;
  execute(interaction: ChatInputCommandInteraction): Promise<void>;
}
