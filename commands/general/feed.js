const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");
const { addSubCommand, subCommandHandling } = require("../../lib/commandUtils");
const { GetLogChannel, GetErrorChannel } = require("../../lib/channelUtils");

module.exports = {
  guildOnly: true,
  handlingCommands: subCommandHandling("admin/feed"),
  data: addSubCommand(
    "admin/feed",
    new SlashCommandBuilder()
      .setName("feed")
      .setDescription("RSS feed/atom feed Utilities")
  ),
  async execute(interaction) {
    const command = this.handlingCommands.get(
      interaction.options.getSubcommand()
    );
    await interaction.deferReply();
    if (!command) {
      console.log(`[-] Not Found: ${interaction.options.getSubcommand()}`);
      return;
    }
    try {
      await command.execute(interaction);
      console.log(`[Run] ${interaction.options.getSubcommand()}`);

      const logEmbed = new EmbedBuilder()
        .setTitle("サブコマンド実行ログ")
        .setDescription(`${interaction.user} がサブコマンドを実行しました。`)
        .setColor(interaction.client.conf.color.s)
        .setTimestamp()
        .setThumbnail(interaction.user.displayAvatarURL({ dynamic: true }))
        .addFields([
          {
            name: "サブコマンド",
            value: `\`\`\`\n/${interaction.options.getSubcommand()}\n\`\`\``,
          },
          {
            name: "実行サーバー",
            value:
              "```\n" + interaction.inGuild()
                ? `${interaction.guild.name} (${interaction.guild.id})`
                : "DM" + "\n```",
          },
          {
            name: "実行ユーザー",
            value:
              "```\n" +
              `${interaction.user.tag}(${interaction.user.id})` +
              "\n```",
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
        .setColor(interaction.conf.color.e)
        .setTimestamp();

      await interaction.editReply(messageEmbed);
      const logChannel = await GetLogChannel(interaction);
      if (logChannel) {
        logChannel.send({ embeds: [messageEmbed] });
      }
    }
    return;
  },
};
