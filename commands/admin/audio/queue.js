const { SlashCommandSubcommandGroupBuilder } = require("discord.js");

const getCmd = require("./queue/get");

module.exports = {
    data: new SlashCommandSubcommandGroupBuilder()
        .setName("queue")
        .setDescription("キューに関するコマンド")
        .addSubcommand(getCmd.data),
    async execute(i, client) {
        const subCommandName = i.options.getSubcommand();
        const command = require("./queue/" + subCommandName);
        command.execute(i, client);
    }
}