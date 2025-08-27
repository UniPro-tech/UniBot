import { MessageFlags, ModalSubmitInteraction } from "discord.js";

export const name = "schedule_create";

export const execute = async (interaction: ModalSubmitInteraction) => {
  const message = interaction.fields.getTextInputValue("message");
  const time = interaction.fields.getTextInputValue("time");

  const timeRegex = /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/;
  if (!timeRegex.test(time)) {
    await interaction.reply({
      content: "時間の形式が正しくありません。`YYYY-MM-DD HH:MM`の形式で入力してください。",
      ephemeral: true,
    });
    return;
  }

  const scheduledTime = new Date(`${time}:00+09:00`);
  if (isNaN(scheduledTime.getTime())) {
    await interaction.reply({
      content: "無効な日時です。もう一度確認してください。",
      ephemeral: true,
    });
    return;
  }

  if (scheduledTime <= new Date()) {
    await interaction.reply({
      content: "過去の日時は指定できません。未来の日時を入力してください。",
      ephemeral: true,
    });
    return;
  }

  const job = await interaction.client.agenda.schedule(scheduledTime, "send discord message", {
    channelId: interaction.channelId,
    message: message,
  });

  await interaction.reply({
    content: `メッセージを${interaction.client.functions.timeUtils.timeToJSTstamp(
      scheduledTime.getTime(),
      true
    )}に送信するようにスケジュールしました。(ジョブID: ${job.attrs._id})`,
    flags: [MessageFlags.Ephemeral],
  });
};
