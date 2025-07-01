import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import config from "@/config";
import { ChatInputCommandInteraction } from "discord.js";

export const guildOnly = false;

export const data = new SlashCommandBuilder()
  .setName("avatar")
  .setDescription("アイコンのURLを取得")
  .addUserOption((option) =>
    option
      .setName("target")
      .setDescription("ここにユーザーを指定してそのユーザーのアイコンのURLを取得します。")
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const user = interaction.options.getUser("target") || interaction.user;
  const embed = new EmbedBuilder()
    .setTitle(`${user.id === interaction.user.id ? "Your" : `${user.username}'s`} Avatar`)
    .setDescription(`URL: ${user.displayAvatarURL()}`)
    .setImage(user.displayAvatarURL())
    .setColor(config.color.success)
    .setTimestamp();

  await interaction.reply({ embeds: [embed] });
  return "No data";
};

export default {
  guildOnly,
  data,
  execute,
};
