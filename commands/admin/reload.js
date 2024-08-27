const { SlashCommandBuilder } = require("discord.js");
const conf = require("../../config.js");
const { addCmd, handling } = require("../../lib/commandUtils.js");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("reload")
    .setDescription("Reloads a command."),
  adminGuildOnly: true,
  /**
   * Executes the reload command.
   *
   * @param {Interaction} i - The interaction object.
   * @returns {Promise<void>} - A promise that resolves after the reload is complete.
   */
  async execute(i) {
    if (!i.member.roles.cache.has(conf.adminRoleId)) {
      i.reply("権限がありません");
      return;
    }
    await i.reply("Now reloading...");
    await addCmd(i.client);
    await handling(i.client);
    await i.editReply("Reloaded!");
  },
};
