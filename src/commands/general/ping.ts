import { SlashCommandBuilder, EmbedBuilder } from "@discordjs/builders";
import { CommandInteraction } from "discord.js";

export const guildOnly = false;
export const name = "ping";
export const execute = async (interaction: CommandInteraction) => {
  const cmdPing = new Date().valueOf() - interaction.createdAt.valueOf();
  const embed = new EmbedBuilder()
    .setTitle("Ping")
    .setDescription("Pong!")
    .addFields([
      {
        name: "WebSocket",
        value: ` ** ${interaction.client.ws.ping} ms ** `,
        inline: true,
      },
      { name: "コマンド受信", value: `** ${cmdPing} ms ** `, inline: true },
    ])
    .setColor(interaction.client.config.color.s)
    .setTimestamp();
  interaction.reply({ embeds: [embed] });
  return "No data";
};
export const data = new SlashCommandBuilder().setName(name).setDescription("Ping値を測定");

export default {
  guildOnly,
  name,
  data,
  execute,
};
