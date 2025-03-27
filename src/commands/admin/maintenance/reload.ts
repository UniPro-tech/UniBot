import { CommandInteraction, GuildMemberRoleManager } from "discord.js";

import { SlashCommandSubcommandBuilder } from "@discordjs/builders";
import conf from "@/config";
import { addCommand, handling } from "@/lib/commandUtils";

export const data = new SlashCommandSubcommandBuilder()
  .setName("reload")
  .setDescription("Reloads a command.");
export const execute = async (interaction: CommandInteraction) => {
  if (
    !(interaction.member?.roles as GuildMemberRoleManager).cache.has(
      conf.adminRoleId
    )
  ) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("Now reloading...");
  await addCommand(interaction.client);
  await handling(interaction.client);
  await interaction.editReply("Reloaded!");
};
