const { SlashCommandBuilder, EmbedBuilder, CommandInteraction } = require("discord.js");
const { subCommandHandling, addSubCommand } = require(`../../lib/commandUtils`);
const { GetLogChannel, GetErrorChannel } = require(`../../lib/channelUtils`);
const config = require("../../config");

module.exports = {
  handlingCommands: subCommandHandling("admin/maintenance"),
  data: addSubCommand(
    "admin/maintenance",
    new SlashCommandBuilder()
      .setName("maintenance")
      .setDescription("メンテナンスモード")
  ),
  adminGuildOnly: true,
  /**
   * Executes the maintenance command.
   *
   * @param {CommandInteraction} interaction - The interaction object.
   * @returns {Promise<void>} - A promise that resolves when the execution is complete.
   */
  async execute(interaction) {
    const command = this.handlingCommands.get(
      interaction.options.getSubcommand()
    );
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

      await interaction.reply(messageEmbed);
      const logChannel = await GetLogChannel(interaction);
      if (logChannel) {
        logChannel.send({ embeds: [messageEmbed] });
      }
    }
    return;
  },
};
