const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("reload")
    .setDescription("Reloads a command."),
  async execute(i, client) {
    if (!i.member.roles.cache.has(client.conf.adminRoleId)) {
      i.reply("権限がありません");
      return;
    }
    i.reply("Now rebooting...");
    process.exit(1);
  },
};
