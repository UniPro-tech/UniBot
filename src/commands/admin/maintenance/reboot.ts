import {
  CommandInteraction,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";

export const data = new SlashCommandSubcommandBuilder()
  .setName("reboot")
  .setDescription("Reboot.");
export const adminGuildOnly = true;
export const execute = async (interaction: CommandInteraction) => {
  if (
    !(interaction.member?.roles as GuildMemberRoleManager).cache.has(
      config.adminRoleId
    )
  ) {
    interaction.reply("権限がありません");
    return;
  }
  await interaction.reply("Now rebooting...");
  process.exit(1);
};

export default {
  data,
  adminGuildOnly,
  execute,
};
