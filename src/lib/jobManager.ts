import { Client, DMChannel, TextChannel, ThreadChannel, VoiceChannel } from "discord.js";
import { agenda } from "..";
import { Job } from "@hokify/agenda";

export interface DiscordMessageJobData {
  channelId: string;
  message: string;
}

export interface RssFeedJobData {
  name: string;
  feedUrl: string;
  channelId: string;
}

export const defineRemindJob = (id: string, client: Client) => {
  agenda.define(`send-discord-message id:${id}`, async (job: Job) => {
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
  await agenda.schedule(when, `send-discord-message id:${id}`, data);
};

export const cronRemindJob = async (id: string, data: DiscordMessageJobData, interval: string) => {
  await agenda.every(interval, `send-discord-message id:${id}`, data);
};

export const cancelRemindJob = async (id: string) => {
  await agenda.cancel({ name: `send-discord-message id:${id}` });
};

export const redefineJobs = async (client: Client) => {
  agenda.jobs({ name: /send-discord-message id:.*/ }).then((jobs) => {
    jobs.forEach((job) => {
      const name = job.attrs.name;
      const id = name.split("id:")[1];
      defineRemindJob(id, client);
    });
  });
};

import { XMLParser } from "fast-xml-parser";

export const defineRssJob = (id: string, client: Client) => {
  agenda.define(`rss-feed id:${id}`, async (job: Job) => {
    const { feedUrl, channelId, name } = job.attrs.data as RssFeedJobData;
    const lastFinishedAt = job.attrs.lastFinishedAt;
    if (!lastFinishedAt) {
      // 初回時は何もしない
      return;
    }
    const channel = client.channels.cache.get(channelId) as
      | TextChannel
      | VoiceChannel
      | ThreadChannel
      | DMChannel;
    if (channel) {
      const response = await fetch(feedUrl);
      const text = await response.text();
      const parser = new XMLParser();
      const feed = parser.parse(text);
      let items = [];
      if (feed.rss && feed.rss.channel && feed.rss.channel.item) {
        items = feed.rss.channel.item;
      } else if (feed.feed && feed.feed.entry) {
        items = feed.feed.entry;
      }
      if (items.length > 0) {
        items.forEach(async (item: any) => {
          if (typeof item === "string") {
            throw new Error("Invalid RSS item format");
          }
          let pubDateStr = "";
          if (item.pubDate) {
            pubDateStr = item.pubDate;
          } else if (item.published) {
            pubDateStr = item.published;
          } else if (item.updated) {
            pubDateStr = item.updated;
          }
          const pubDate = new Date(pubDateStr);
          if (pubDate <= lastFinishedAt) {
            return;
          }
          let title = "";
          let link = "";
          if (item.title) {
            title = typeof item.title === "string" ? item.title : item.title["#text"];
          }
          if (item.link) {
            link =
              typeof item.link === "string" ? item.link : item.link["@_href"] || item.link["#text"];
          }
          await channel.send(`${name}に新しい記事があります: **${title}**\n${link}`);
        });
      }
    }
  });
};

export const redefineJobsRss = async (client: Client) => {
  agenda.jobs({ name: /rss-feed id:.*/ }).then((jobs) => {
    jobs.forEach((job) => {
      const name = job.attrs.name;
      const id = name.split("id:")[1];
      defineRssJob(id, client);
    });
  });
};

export const addRssJob = async (id: string, data: RssFeedJobData, interval: string) => {
  await agenda.jobs({ name: `rss-feed id:${id}` }).then(async (jobs) => {
    for (const job of jobs) {
      await job.run();
    }
  });
  await agenda.every(interval, `rss-feed id:${id}`, data);
};

export default {
  defineRemindJob,
  scheduleRemindJob,
  cronRemindJob,
  cancelRemindJob,
  addRssJob,
  redefineJobsRss,
};
