import { addRssJob, defineRssJob } from "@/lib/jobManager";
import {
  ChatInputCommandInteraction,
  SlashCommandSubcommandBuilder,
  TextChannel,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("set")
  .setDescription("RSSフィードの追加")
  .addStringOption((option) =>
    option.setName("url").setDescription("RSSフィードのURL").setRequired(true)
  )
  .addStringOption((option) =>
    option.setName("name").setDescription("フィードの名前").setRequired(true)
  )
  .addChannelOption((option) =>
    option
      .setName("channel")
      .setDescription("RSSフィードを投稿するチャンネル(デフォルトは現在のチャンネル)")
      .setRequired(false)
  )
  .addStringOption((option) =>
    option.setName("interval").setDescription("フィードの更新間隔（分）").setRequired(false)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const feedUrl = interaction.options.getString("url", true);
  const name = interaction.options.getString("name", true);
  const channel = interaction.options.getChannel("channel") || interaction.channel;
  const interval = `*/${interaction.options.getString("interval") || 30} * * * *`;

  if (!channel || channel instanceof TextChannel === false) {
    await interaction.reply({
      content: "有効なテキストチャンネルを指定してください。",
      ephemeral: true,
    });
    return;
  }

  await defineRssJob(interaction.id, interaction.client);
  await addRssJob(interaction.id, { feedUrl, channelId: channel.id, name }, interval);

  await interaction.reply({
    content: `RSSフィード「${name}」がチャンネル ${channel} に追加されました。更新間隔: ${interval}`,
    ephemeral: true,
  });
};
