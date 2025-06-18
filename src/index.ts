//import cron from "node-cron";
import {
  Client,
  GatewayIntentBits,
  Collection,
  Partials,
  EmbedBuilder,
  Channel,
  TextChannel,
} from "discord.js";

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildVoiceStates,
  ],
  partials: [Partials.Channel],
});

import fs from "fs";
import config from "@/config";
import timeUtils from "@/lib/timeUtils";
import logUtils from "@/lib/dataUtils";
client.config = config;
client.functions = {
  timeUtils: timeUtils,
  logUtils: logUtils,
};
client.fs = fs;

import {
  ChatInputCommandCollector,
  MessageContextMenuCommandCollector,
  StringSelectMenuCollector,
} from "@/lib/collecter";
import path from "path";
import { ChatInputCommand } from "./executors/types/ChatInputCommand";
import { StringSelectMenu } from "./executors/types/StringSelectMenu";
client.interactionExecutorsCollections = {
  chatInputCommands: new Collection<string, ChatInputCommand>(),
  stringSelectMenus: new Collection<string, StringSelectMenu>(),
  // TODO: ここはMessageContextMenuCommandにする
  messageContextMenuCommands: new Collection<string, ChatInputCommand>(),
};
ChatInputCommandCollector(client);
StringSelectMenuCollector(client);
MessageContextMenuCommandCollector(client);
const eventFiles = fs
  .readdirSync(path.resolve(__dirname, "events"))
  .filter((file) => (file.endsWith(".ts") && !file.endsWith(".d.ts")) || file.endsWith(".js"));
for (const file of eventFiles) {
  const event = require(path.resolve(__dirname, `./events/${file}`));
  if (event.once) {
    try {
      client.once(event.name, (...args) => event.execute(...args, client));
    } catch (error) {
      console.error(
        `\u001b[31m[${client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}]\u001b[0m\n`,
        error
      );
    }
  } else {
    try {
      client.on(event.name, (...args) => event.execute(...args, client));
    } catch (error) {
      console.error(
        `\u001b[31m[${client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}]\u001b[0m\n`,
        error
      );
    }
  }
}

// TODO:ここ下の3行のコメントアウトを外し、いい感じにする
//const { rssGet } = require("./lib/rss.cjs");

//cron.schedule("*/1-59 * * * *", async () => {
/*  console.log("Cron job start");
  const files = fs.readdirSync("./log/v1/feed");
  for (const file of files) {
    const datas = await JSON.parse(
      fs.readFileSync(`./log/v1/feed/${file}.log`).toString()
    );
    datas.forEach(async (data, index) => {
      console.log(data.url);
      const items = await rssGet(data.url);
      const channel = client.channels.cache.get(file.replace(".log", ""));
      for (const item of items) {
        if (item.pubDate <= data.lastDate) continue;
        console.log("send");
        const embed = new EmbedBuilder()
          .setTitle(item.title)
          .setURL(item.link)
          .setDescription(item.content)
          .setColor(config.color.success)
          .setTimestamp();
        channel.send({ embeds: [embed] });
        datas[index].lastDate = items[0].pubDate;
        client.function.logUtils.write(datas, `v1/feed/${file}`);
      }
    });
  }
});
*/

// エラー処理 (これ入れないとエラーで落ちる。本当は良くないかもしれない)
process.on("uncaughtException", (error) => {
  console.error(`[${client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}] ${error.stack}`);
  const embed = new EmbedBuilder()
    .setTitle("ERROR - uncaughtException")
    .setDescription("```\n" + error.stack + "\n```")
    .setColor(config.color.error)
    .setTimestamp();
  client.channels.fetch(config.logch.error).then((channel: Channel | null) => {
    if (!channel || !(channel instanceof TextChannel)) {
      console.error("Error: Log Channel is invalid.");
      return;
    }
    channel.send({ embeds: [embed] });
  });
});

process.on("unhandledRejection", (reason, promise) => {
  console.error(
    `\u001b[31m[${client.functions.timeUtils.timeToJSTstamp(
      Date.now(),
      true
    )}] ${reason}\u001b[0m\n`,
    promise
  );
  const embed = new EmbedBuilder()
    .setTitle("ERROR - unhandledRejection")
    .setDescription("```\n" + reason + "\n```")
    .setColor(config.color.error)
    .setTimestamp();
  client.channels.fetch(config.logch.error).then((channel: Channel | null) => {
    if (!channel || !(channel instanceof TextChannel)) {
      console.error("Error: Log Channel is invalid.");
      return;
    }
    channel.send({ embeds: [embed] });
  });
});

client.login(config.token);
