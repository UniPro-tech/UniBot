const { SlashCommandSubcommandBuilder } = require("discord.js");
const conf = require(`../../../config.js`);

module.exports = {
  data: new SlashCommandSubcommandBuilder()
    .setName("reboot")
    .setDescription("Reboot."),
  adminGuildOnly: true,
  async execute(i) {
    if (!i.member.roles.cache.has(conf.adminRoleId)) {
      i.reply("権限がありません");
      return;
    }
    await i.reply("Now rebooting...");
    process.exit(1);
  },
};
