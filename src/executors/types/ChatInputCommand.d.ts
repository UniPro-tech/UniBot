import {
  Collection,
  CommandInteraction,
  SlashCommandBuilder,
  SlashCommandSubcommandBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";

export interface ChatInputCommand {
  handlingCommands: Collection<string, ChatInputCommand> | undefined;
  adminGuildOnly: boolean | undefined;
  guildOnly: boolean | undefined;
  data: SlashCommandBuilder | SlashCommandSubcommandBuilder | SlashCommandSubcommandGroupBuilder;
  execute(interaction: CommandInteraction): Promise<void>;
}
