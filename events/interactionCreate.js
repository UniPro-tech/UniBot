const { Events, EmbedBuilder } = require("discord.js");
const config = require("../config");
const { GetLogChannel, GetErrorChannel } = require(`../lib/channelUtils`);

module.exports = {
    name: Events.InteractionCreate,
    async execute(interaction) {
        if (interaction.isChatInputCommand()) {
            console.log(`[!] ${interaction.commandName}`);
            const command = interaction.client.commands.get(interaction.commandName);
            if (!command) {
                console.log(`[-] Not Found: ${interaction.commandName}`);
                return;
            }
            if (!interaction.inGuild() && command.guildOnly) {
                const embed = new EmbedBuilder()
                    .setTitle("エラー")
                    .setDescription("このコマンドはDMでは実行できません。")
                    .setColor(interaction.client.config.color.e);
                interaction.reply({ embeds: [embed] });
                return;
            }

            try {
                await command.execute(interaction);
                console.log(`[Run] ${interaction.commandName}`);

                const logEmbed = new EmbedBuilder()
                    .setTitle("コマンド実行ログ")
                    .setDescription(`${interaction.user} がコマンドを実行しました。`)
                    .setColor(config.color.s)
                    .setTimestamp()
                    .setThumbnail(interaction.user.displayAvatarURL({ dynamic: true }))
                    .addFields([
                        {
                            name: "コマンド",
                            value: `\`\`\`\n/${interaction.commandName}\n\`\`\``,
                        },
                        {
                            name: "実行サーバー",
                            value: "```\n"
                                + interaction.inGuild() ? `${interaction.guild.name} (${interaction.guild.id})` : "DM"
                            + "\n```",
                        },
                        {
                            name: "実行ユーザー",
                            value: "```\n" + `${interaction.user.tag}(${interaction.user.id})` + "\n```",
                        },
                    ])
                    .setFooter({ text: `${interaction.id}` });
                const channel = await GetLogChannel(interaction);
                if (channel) {
                    channel.send({ embeds: [logEmbed] });
                }
            } catch (error) {
                console.error(error);
                const logEmbed = new EmbedBuilder()
                    .setTitle("ERROR - cmd")
                    .setDescription("```\n" + error.toString() + "\n```")
                    .setColor(config.color.e)
                    .setTimestamp();

                const channel = await GetErrorChannel(interaction);
                if (channel) {
                    channel.send({ embeds: [logEmbed] });
                }
                const messageEmbed = new EmbedBuilder()
                    .setTitle("すみません、エラーが発生しました...")
                    .setDescription("```\n" + error + "\n```")
                    .setColor(config.color.e)
                    .setTimestamp();

                await interaction.reply(messageEmbed);
                const logChannel = await GetLogChannel(interaction);
                if (logChannel) {
                    logChannel.send({ embeds: [messageEmbed] });
                }
            }
        }
    }
};