import {
  CommandInteraction,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";

import config from "@/config";
import { addCommand, handling } from "@/lib/commandUtils";

export const data = new SlashCommandSubcommandBuilder()
  .setName("reload")
  .setDescription("Reloads a command.");
export const adminGuildOnly = true;
export const execute = async (interaction: CommandInteraction) => {
  if (!(interaction.member?.roles as GuildMemberRoleManager).cache.has(config.adminRoleId)) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("Now reloading...");
  await addCommand(interaction.client);
  await handling(interaction.client);
  await interaction.editReply("Reloaded!");
};

export default {
  data,
  adminGuildOnly,
  execute,
};
