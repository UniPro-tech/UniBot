import { MessageFlags, ModalSubmitInteraction } from "discord.js";
import later from "@breejs/later";
import { ALStorage, loggingSystem } from "@/index";

export const name = "schedule_create_repeat";

export const execute = async (interaction: ModalSubmitInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "modals/schedule_create_repeat" });
  const message = interaction.fields.getTextInputValue("message");
  const time = interaction.fields.getTextInputValue("time");
  const cronText = convertToCron(time);
  if (!cronText) {
    await interaction.reply({
      content: "時間の形式が不正です。もう一度やり直してください。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const job = await interaction.client.functions.jobManager.defineRemindJob(
    interaction.id,
    interaction.client
  );
  const data = {
    channelId: interaction.channelId as string,
    message: message,
  };
  await interaction.client.functions.jobManager.cronRemindJob(interaction.id, data, cronText);

  logger.info(
    { service: "ScheduleCreateRepeat", jobId: interaction.id, cron: cronText, message },
    "Scheduled repeating message"
  );

  await interaction.reply({
    content: `メッセージを${time}に送信するようにスケジュールしました。(ジョブID: ${interaction.id})`,
    flags: [MessageFlags.Ephemeral],
  });
};

function convertToCron(laterText: string): string | null {
  // "every day"だけが含まれてて他に何もなければ9:00 amに変換する
  if (/^\s*every day\s*$/i.test(laterText.trim())) {
    laterText = "every day at 9:00 am";
  }
  // every dayが含まれているが他に何かしらの指定もある場合、every dayを削除する
  if (/every day/i.test(laterText.trim()) && !/^\s*every day\s*$/i.test(laterText.trim())) {
    laterText = laterText.replace(/every day/gi, "").trim();
  }
  // 0時台（例: 0:00, 0:15, 0:59 など）があれば12時台(am)に変換し、am/pmがなければamをつける
  laterText = laterText.replace(/\b0:([0-5][0-9])\b/g, "12:$1");
  // 12:xx でam/pmがついてないやつにamをつける
  laterText = laterText.replace(/\b12:([0-5][0-9])\b(?!\s?(am|pm))/gi, "12:$1 am");
  // 13時以降の時刻が含まれていたらam/pm記法に変換する
  laterText = laterText.replace(/\b([1][3-9]|2[0-3]):([0-5][0-9])\b/g, (_match, hour, minute) => {
    const h = parseInt(hour, 10);
    const ampmHour = h > 12 ? h - 12 : h;
    const period = h >= 12 ? "pm" : "am";
    return `${ampmHour}:${minute} ${period}`;
  });
  const sched = later.parse.text(laterText);
  if (!sched) {
    return null;
  }
  if (sched.error > -1) {
    return null;
  }

  const nextTwo = later.schedule(sched).next(2) as Date[];
  if (nextTwo.length < 2) {
    return null;
  }

  const first = nextTwo[0];
  const second = nextTwo[1];
  const diffMs = second.getTime() - first.getTime();
  const diffMinutes = Math.round(diffMs / 60000);

  const min = first.getMinutes();
  const hour = first.getHours();
  const date = first.getDate();
  const month = first.getMonth() + 1; // cronは1-12
  const weekDay = first.getDay();

  // 1分〜59分ごと
  if (diffMinutes < 60) {
    return `*/${diffMinutes} * * * *`;
  }
  // 1時間〜23時間ごと
  if (diffMinutes % 60 === 0 && diffMinutes < 1440) {
    const hours = diffMinutes / 60;
    return `${min} */${hours} * * *`;
  }
  // 毎日同時刻
  if (diffMinutes >= 1440 && diffMinutes < 10080) {
    return `${min} ${hour} * * *`;
  }
  // 毎週同時刻
  if (diffMinutes >= 10080 && diffMinutes < 40320) {
    return `${min} ${hour} * * ${weekDay}`;
  }
  // 毎月同日同時刻
  if (diffMinutes >= 40320 && diffMinutes < 525600) {
    return `${min} ${hour} ${date} * *`;
  }
  // 毎年同月同日同時刻
  if (diffMinutes >= 525600) {
    return `${min} ${hour} ${date} ${month} *`;
  }

  return null;
}
