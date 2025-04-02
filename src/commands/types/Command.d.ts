import {
  Collection,
  CommandInteraction,
  SlashCommandBuilder,
  SlashCommandSubcommandBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";

export interface Command {
  handlingCommands: Collection<string, Command> | undefined;
  adminGuildOnly: boolean | undefined;
  guildOnly: boolean | undefined;
  data:
    | SlashCommandBuilder
    | SlashCommandSubcommandBuilder
    | SlashCommandSubcommandGroupBuilder;
  execute(interaction: CommandInteraction): Promise<void>;
}
