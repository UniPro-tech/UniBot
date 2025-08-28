import {
  ChatInputCommandInteraction,
  EmbedBuilder,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";

export const data = new SlashCommandSubcommandBuilder()
  .setName("list")
  .setDescription("予約投稿の一覧を表示");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  await interaction.deferReply({ ephemeral: true });
  interaction.client.agenda.now("purge agenda");
  await new Promise((resolve) => setTimeout(resolve, 5000));
  const jobs = await interaction.client.agenda.jobs({
    name: /send-discord-message id:.*/,
  });

  if (!jobs.length) {
    await interaction.editReply({
      content: "予約投稿はまだありません。",
    });
    return;
  }

  const embed = new EmbedBuilder()
    .setTitle("予約投稿一覧")
    .setColor(config.color.success)
    .setTimestamp();

  jobs.forEach((job) => {
    embed.addFields({
      name: `ジョブID: ${job.attrs.name?.toString().trim().split("id:")[1] || "不明"}`,
      value: `内容: ${
        (job.attrs.data as { channelId: string; message: string })?.message || "不明"
      }\n投稿予定: <t:${Math.floor((job.attrs.nextRunAt?.getTime() ?? 0) / 1000)}:F>`,
    });
  });

  await interaction.editReply({
    embeds: [embed],
  });
};

export default {
  data,
  execute,
};
