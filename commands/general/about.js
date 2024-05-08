const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("about")
    .setDescription("このBotについての情報を表示します。"),
  async execute(interaction) {
    const size = await interaction.client.shard.fetchClientValues("guilds.cache.size");
    const embed = new EmbedBuilder()
      .setColor(0x0099ff)
      .setTitle(`About ${interaction.client.conf.productname}`)
      .setURL("https://uniproject-tech.net/UniBot/")
      .setAuthor(interaction.client.conf.author)
      .setDescription(interaction.client.conf.description)
      .setThumbnail(interaction.client.user.displayAvatarURL({ dynamic: true }))
      .addFields(
        { name: "Version", value: interaction.client.conf.version },
        { name: "Author", value: interaction.client.conf.author.name },
        { name: "Guild Count", value: `${size}` }
      )
      .setTimestamp();
    interaction.reply({ embeds: [embed] });
  },
};
