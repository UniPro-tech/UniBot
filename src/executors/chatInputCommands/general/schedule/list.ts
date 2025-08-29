import {
  ChatInputCommandInteraction,
  EmbedBuilder,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";
import cronstrue from "cronstrue";

export const data = new SlashCommandSubcommandBuilder()
  .setName("list")
  .setDescription("予約投稿の一覧を表示");

export const execute = async (interaction: ChatInputCommandInteraction) => {
  await interaction.deferReply({ ephemeral: true });
  interaction.client.agenda.now("purge agenda");
  await new Promise((resolve) => setTimeout(resolve, 2000));
  const jobs = await interaction.client.agenda.jobs({
    name: /send-discord-message id:.*/,
    "data.channelId": interaction.channel?.id.toString(),
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
    const jobId = job.attrs.name?.toString().trim().split("id:")[1] || "不明";
    const message = (job.attrs.data as { channelId: string; message: string })?.message || "不明";
    const nextRunAt = job.attrs.nextRunAt?.getTime() || 0;
    const cronInterval = job.attrs.repeatInterval?.toString() || null;
    embed.addFields({
      name: `ジョブID: ${jobId}`,
      value: `メッセージ: ${message}\n次回実行予定: <t:${Math.floor(
        nextRunAt / 1000
      )}:F>\n繰り返し: ${cronInterval ? cronstrue.toString(cronInterval) : "いいえ"}`,
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
