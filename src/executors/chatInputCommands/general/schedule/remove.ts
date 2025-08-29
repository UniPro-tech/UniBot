import {
  ChatInputCommandInteraction,
  EmbedBuilder,
  MessageFlags,
  SlashCommandSubcommandBuilder,
} from "discord.js";
import config from "@/config";

export const data = new SlashCommandSubcommandBuilder()
  .setName("remove")
  .setDescription("予約投稿を削除")
  .addStringOption((option) =>
    option.setName("jobid").setDescription("削除する予約投稿のジョブID").setRequired(true)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  const jobId = interaction.options.getString("jobid");
  if (!jobId) {
    await replyWithError(
      interaction,
      "Error - ジョブID未指定",
      "削除する予約投稿のジョブIDを指定してください。"
    );
    return;
  }

  await interaction.client.functions.jobManager.cancelRemindJob(jobId);

  await interaction.reply({
    content: `ジョブID「${jobId}」の予約投稿を削除しました。`,
    flags: [MessageFlags.Ephemeral],
  });
};

export default {
  data,
  execute,
};

const replyWithError = async (
  interaction: ChatInputCommandInteraction,
  title: string,
  description: string
) => {
  const embed = new EmbedBuilder()
    .setTitle(title)
    .setDescription(description)
    .setColor(config.color.error)
    .setTimestamp();
  await interaction.reply({
    embeds: [embed],
    flags: [MessageFlags.Ephemeral],
  });
};
