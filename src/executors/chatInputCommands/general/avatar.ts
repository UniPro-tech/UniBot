import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import config from "@/config";
import { ChatInputCommandInteraction, CommandInteractionOptionResolver } from "discord.js";

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
  const user = (interaction.options as CommandInteractionOptionResolver).getUser("target");
  if (user) {
    const embed = new EmbedBuilder()
      .setTitle(`${user.username}'s Avatar`)
      .setDescription(`URL:${user.displayAvatarURL()}`)
      .setImage(user.displayAvatarURL())
      .setColor(config.color.success)
      .setTimestamp();
    interaction.reply({ embeds: [embed] });
  } else {
    const embed = new EmbedBuilder()
      .setTitle(`Your Avatar`)
      .setDescription(`URL:${interaction.user.displayAvatarURL()}`)
      .setImage(`${interaction.user.displayAvatarURL()}`)
      .setColor(config.color.success)
      .setTimestamp();
    interaction.reply({ embeds: [embed] });
  }
  return "No data";
};

export default {
  guildOnly,
  data,
  execute,
};
