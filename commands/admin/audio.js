const { SlashCommandBuilder } = require("discord.js");

const playCmd = require("./audio/play");
const skipCmd = require("./audio/skip");
const queueGroup = require("./audio/queue");

module.exports = {
    guildOnly: true, // サーバー専用コマンドかどうか
    adminGuildOnly: true,
    data: new SlashCommandBuilder()
        .setName("audio")
        .setDescription(
            "オーディオに関連するコマンド"
        )
        .addSubcommand(playCmd.data)
        .addSubcommand(skipCmd.data)
        .addSubcommandGroup(queueGroup.data),
    async execute(i, client) {
        const subCommandGroupName = i.options.getSubcommandGroup();
        if (!subCommandGroupName) {
            const subCommandName = i.options.getSubcommand();
            const command = require("./audio/" + subCommandName);
            command.execute(i, client);
        } else {
            const command = require("./audio/" + subCommandGroupName);
            command.execute(i, client);
        }
    }
}