import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";

export const guildOnly = false;
export const name = "ping";

export const execute = async (interaction: CommandInteraction) => {
  const cmdPing = Date.now() - interaction.createdAt.valueOf();
  const wsPing = interaction.client.ws.ping;
  const color = interaction.client.config.color.success;

  const embed = new EmbedBuilder()
    .setTitle("Ping")
    .setDescription("Pong!")
    .addFields(
      { name: "WebSocket", value: `**${wsPing} ms**`, inline: true },
      { name: "コマンド受信", value: `**${cmdPing} ms**`, inline: true }
    )
    .setColor(color)
    .setTimestamp();

  await interaction.reply({ embeds: [embed] });
  return "No data";
};

export const data = new SlashCommandBuilder().setName(name).setDescription("Ping!");

export default {
  guildOnly,
  name,
  data,
  execute,
};
