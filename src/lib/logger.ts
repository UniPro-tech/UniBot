import { ChatInputCommandInteraction, Client, EmbedBuilder } from "discord.js";
import { GetLogChannel } from "./channelUtils";

export const writeChatInputCommandInteractionLog = async (
  interaction: ChatInputCommandInteraction
) => {
  const userId = interaction.user.id;
  const usertag = interaction.user.tag;
  const commandName = interaction.commandName;
  const channel = interaction.channel?.id;
  const channelName =
    interaction.channel && interaction.channel.type === 1
      ? "DM"
      : "name" in (interaction.channel ?? {}) &&
        typeof (interaction.channel as any).name === "string"
      ? (interaction.channel as { name: string }).name
      : "Unknown";
  const timestamp = new Date();
  const subcommand = interaction.options.getSubcommand(false);
  const subcommandGroup = interaction.options.getSubcommandGroup(false);
  const guild = interaction.guild?.id;
  const guildName = interaction.guild?.name;
  const interactionId = interaction.id;
  const userIcon = interaction.user.displayAvatarURL();
  const options = interaction.options.data
    .filter(
      (option) =>
        option.value !== null &&
        option.value !== undefined &&
        option.name !== "subcommand" &&
        option.name !== "subcommandGroup"
    )
    .sort((a, b) => a.name.localeCompare(b.name))
    .map((option) => ({
      name: option.name,
      value: `\`\`\`\n${option.value}\n\`\`\``,
    }));
  const embed = new EmbedBuilder()
    .setColor(interaction.client.config.color.success)
    .setTitle("コマンド実行ログ")
    .setThumbnail(userIcon)
    .setDescription(`<@${userId}> が <#${channel}> でコマンドを実行しました`)
    .addFields([
      {
        name: "User",
        value: `\`\`\`\n${usertag} (${userId})\n\`\`\``,
      },
      {
        name: "Command",
        value: `\`\`\`\n\/${commandName} ${subcommandGroup || ""} ${
          subcommand || ""
        }\n\`\`\``,
      },
      {
        name: "Channel",
        value: `\`\`\`\n${channelName} (${channel})\n\`\`\``,
      },
      {
        name: "Guild",
        value: guild ? `\`\`\`\n${guildName} (${guild})\n\`\`\`` : "DM",
      },
    ])
    .addFields(options as [{ name: string; value: string }])
    .setTimestamp(timestamp)
    .setFooter({ text: `Interaction ID: ${interactionId}` });
  const logChannel = await GetLogChannel(interaction.client);
  if (logChannel) {
    await logChannel.send({ embeds: [embed] });
  }
};
