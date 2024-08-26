const { SlashCommandBuilder, EmbedBuilder } = require("discord.js");
const package_data = require("../../package.json");

module.exports = {
  data: new SlashCommandBuilder()
    .setName("about")
    .setDescription("このBotについての情報を表示します。"),
  async execute(interaction) {
    //const size = await interaction.client.shard.fetchClientValues("guilds.cache.size");
    const embed = new EmbedBuilder()
      .setColor(0x0099ff)
      .setTitle(`About ${package_data.name}`)
      .setTimestamp();
    if (interaction.client.shard) {
      embed.addFields({name: "Shard ID", value: interaction.client.shard.ids, inline: true});
      embed.addFields({name: "Shard Count", value: interaction.client.shard.count, inline: true});
    }
    if (package_data.description) {
      embed.setDescription(package_data.description);
    }
    if (package_data.version) {
      embed.addFields({ name: "Version", value: package_data.version, inline: true});
    }
    if (package_data.author) {
      embed.addFields({name: "Author", value: `[${package_data.author + package_data.email}](${package_data.author.url})`});
    }
    if (package_data.license) {
      embed.addFields({name: "License", value: package_data.license, inline: true});
    }
    if (package_data.repository) {
      embed.addFields({name: "Repository", value: package_data.repository.url});
    }
    if (package_data.homepage) {
      embed.setURL(package_data.homepage);
    }
    if (package_data.contributors) { 
      let temp = new Array();
      for(let i = 0; i < package_data.contributors.length; i++) {
        temp[i] = `[${ package_data.contributors[i].name} <${package_data.contributors[i].email }>](${ package_data.contributors[i].url })`;
      }
      embed.addFields({name: "Contributors", value: temp.join("\n")});
    }
    interaction.reply({ embeds: [embed] });
  },
};
