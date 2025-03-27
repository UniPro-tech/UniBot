import {
  CommandInteraction,
  GuildMemberRoleManager,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import conf from "@/config";

module.exports = {
  data: new SlashCommandSubcommandBuilder()
    .setName("reboot")
    .setDescription("Reboot."),
  adminGuildOnly: true,
  async execute(interaction: CommandInteraction) {
    if (
      !(interaction.member?.roles as GuildMemberRoleManager).cache.has(
        conf.adminRoleId
      )
    ) {
      interaction.reply("権限がありません");
      return;
    }
    await interaction.reply("Now rebooting...");
    process.exit(1);
  },
};
