import {
  CommandInteraction,
  SlashCommandBuilder,
  SlashCommandSubcommandBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";

export interface ContextMenuCommand {
  adminGuildOnly: boolean | undefined;
  guildOnly: boolean | undefined;
  data: SlashCommandBuilder | SlashCommandSubcommandBuilder | SlashCommandSubcommandGroupBuilder;
  execute(interaction: CommandInteraction): Promise<void>;
}
