import { Client, DMChannel, TextChannel, ThreadChannel, VoiceChannel } from "discord.js";
import { agenda } from "..";
import { Job } from "@hokify/agenda";

export interface DiscordMessageJobData {
  channelId: string;
  message: string;
}

export const defineRemindJob = (id: string, client: Client) => {
  agenda.define(`send discord message id:${id}`, async (job: Job) => {
    const { channelId, message } = job.attrs.data as DiscordMessageJobData;
    const channel = client.channels.cache.get(channelId) as
      | TextChannel
      | VoiceChannel
      | ThreadChannel
      | DMChannel;
    if (channel) {
      await channel.send(message);
    }
  });
};

export const scheduleRemindJob = async (id: string, data: DiscordMessageJobData, when: string) => {
  await agenda.schedule(when, `send discord message id:${id}`, data);
};

export const cronRemindJob = async (id: string, data: DiscordMessageJobData, interval: string) => {
  await agenda.every(interval, `send discord message id:${id}`, data);
};

export const cancelRemindJob = async (id: string) => {
  await agenda.cancel({ name: `send discord message id:${id}` });
};

export default {
  defineRemindJob,
  scheduleRemindJob,
  cronRemindJob,
  cancelRemindJob,
};
