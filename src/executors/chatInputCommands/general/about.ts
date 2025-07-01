import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import { ChatInputCommandInteraction } from "discord.js";
import path from "path";

export const guildOnly = false;

export const data = new SlashCommandBuilder()
  .setName("about")
  .setDescription("Botについての情報を表示");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const packagePath = path.resolve(__dirname, "../../../../package.json");
  const packageData = await import(packagePath);

  const embed = new EmbedBuilder()
    .setColor(0x0099ff)
    .setTitle(`About ${packageData.name}`)
    .setTimestamp()
    .setThumbnail(interaction.client.user?.displayAvatarURL({ size: 1024 }) ?? "")
    .setFooter({
      text: `Requested by ${interaction.user.tag}`,
      iconURL: interaction.user.displayAvatarURL(),
    });

  if (interaction.client.shard) {
    embed.addFields(
      {
        name: "Shard ID",
        value: interaction.client.shard.ids.toString(),
        inline: true,
      },
      {
        name: "Shard Count",
        value: interaction.client.shard.count.toString(),
        inline: true,
      }
    );
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
      value: packageData.author.name ?? packageData.author,
      inline: false,
    });
  }

  if (packageData.repository?.url) {
    embed.addFields({
      name: "Repository",
      value: packageData.repository.url,
      inline: false,
    });
  }

  if (packageData.homepage) {
    embed.setURL(packageData.homepage);
  }

  if (Array.isArray(packageData.contributors) && packageData.contributors.length > 0) {
    const contributors = packageData.contributors.map((c: any) => {
      const name = c.name ?? "";
      const url = c.url ? `[${name}](${c.url})` : name;
      const email = c.email ? `<[${c.email}](mailto:${c.email})>` : "";
      return `- ${url} ${email}`.trim();
    });
    embed.addFields({
      name: "Contributors",
      value: contributors.join("\n"),
      inline: false,
    });
  }

  await interaction.reply({ embeds: [embed] });
};
