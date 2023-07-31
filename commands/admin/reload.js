const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("reload")
    .setDescription("Reloads a command.")
    /*.addStringOption((option) =>
      option
        .setName("command")
        .setDescription("The command to reload.")
        .setRequired(true)
    */,
  async execute(i, client) {
    const add = require(`../../system/add.js`);
    const embed = new EmbedBuilder()
      .setTitle("ReloadCommands!!")
      .setColor(client.conf.color.s)
      .setTimestamp();
    i.reply({ embeds: [embed] });
    add.addCmd(client.conf);
    return null;
  },
};