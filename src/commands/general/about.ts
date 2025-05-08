import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";
import path from "path";
const packageData = require(path.resolve(__dirname, "../../../package.json"));

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
  if (packageData.license) {
    embed.addFields({
      name: "License",
      value: packageData.license,
      inline: true,
    });
  }
  if (packageData.author) {
    embed.addFields({
      name: "Authors",
      value: packageData.author.name,
      inline: false,
    });
  }
  if (packageData.repository) {
    embed.addFields({
      name: "Repository",
      value: packageData.repository.url,
      inline: false,
    });
  }
  if (packageData.homepage) {
    embed.setURL(packageData.homepage);
  }
  if (packageData.contributors) {
    let temp = [];
    for (let i = 0; i < packageData.contributors.length; i++) {
      temp[
        i
      ] = `- [${packageData.contributors[i].name}](${packageData.contributors[i].url}) <[${packageData.contributors[i].email}](mailto:${packageData.contributors[i].email})>`;
    }
    embed.addFields({ name: "Contributors", value: temp.join("\n") });
  }
  embed.setTimestamp();
  embed.setThumbnail(
    interaction.client.user?.displayAvatarURL({
      size: 1024,
    })
  );
  embed.setFooter({
    text: `Requested by ${interaction.user.tag}`,
    iconURL: interaction.user.displayAvatarURL(),
  });
  interaction.reply({ embeds: [embed] });
};
