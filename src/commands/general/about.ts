import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";
const packageData = require("@/package.json");

export const guildOnly = false;
export const data = new SlashCommandBuilder()
  .setName("about")
  .setDescription("このBotについての情報を表示します。");
export const execute = async (interaction: CommandInteraction) => {
  //const size = await interaction.client.shard.fetchClientValues("guilds.cache.size");
  const embed = new EmbedBuilder()
    .setColor(0x0099ff)
    .setTitle(`About ${packageData.name}`)
    .setTimestamp();
  if (interaction.client.shard) {
    embed.addFields([
      {
        name: "Shard ID",
        value: interaction.client.shard.ids.toString(),
        inline: true,
      },
      {
        name: "Shard Count",
        value: interaction.client.shard.count.toString(),
        inline: true,
      },
    ]);
  }
  if (packageData.description) {
    embed.setDescription(packageData.description);
  }
  if (packageData.version) {
    embed.addFields({
      name: "Version",
      value: packageData.version,
      inline: true,
    });
  }
  if (packageData.author) {
    embed.addFields({
      name: "Author",
      value: `[${packageData.author + packageData.email}](${
        packageData.author.url
      })`,
    });
  }
  if (packageData.license) {
    embed.addFields({
      name: "License",
      value: packageData.license,
      inline: true,
    });
  }
  if (packageData.repository) {
    embed.addFields({
      name: "Repository",
      value: packageData.repository.url,
    });
  }
  if (packageData.homepage) {
    embed.setURL(packageData.homepage);
  }
  if (packageData.contributors) {
    let temp = new Array();
    for (let i = 0; i < packageData.contributors.length; i++) {
      temp[
        i
      ] = `[${packageData.contributors[i].name} <${packageData.contributors[i].email}>](${packageData.contributors[i].url})`;
    }
    embed.addFields({ name: "Contributors", value: temp.join("\n") });
  }
  interaction.reply({ embeds: [embed] });
};
