const cron = require("node-cron");

const {
  Client,
  GatewayIntentBits,
  Collection,
  Partials,
  EmbedBuilder,
} = require("discord.js");

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildVoiceStates,
  ],
  partials: [Partials.Channel],
  ws: {
    properties: {
      $os: "Uni OS",
      $browser: "Uni Browser",
      $device: "K8s on Proxmox VE 6.4 in Sibainu",
    },
  },
});

const fs = require("fs");

const config = require("./config.js");
const functions = {
  timeUtils: require("./lib/timeUtils.js"),
  logUtils: require("./lib/logUtils.js"),
};

client.conf = config;
client.func = functions;
client.fs = fs;

const cmdH = require(`./lib/commandUtils.js`);
client.commands = new Collection();
cmdH.handling(client, fs, Collection, config);
const eventFiles = fs
  .readdirSync("./events")
  .filter((file) => file.endsWith(".js"));
for (const file of eventFiles) {
  const event = require(`./events/${file}`);
  if (event.once) {
    try {
      client.once(event.name, (...args) => event.execute(...args, client));
    } catch (error) {
      console.error(
        `\u001b[31m[${functions.timeUtils.timeToJST(
          Date.now(),
          true
        )}]\u001b[0m\n`,
        error
      );
    }
  } else {
    try {
      client.on(event.name, (...args) => event.execute(...args, client));
    } catch (error) {
      console.error(
        `\u001b[31m[${functions.timeUtils.timeToJST(
          Date.now(),
          true
        )}]\u001b[0m\n`,
        error
      );
    }
  }
}

const { rssGet } = require("./lib/rss.cjs");

cron.schedule("*/1-59 * * * *", async () => {
  console.log("Cron job start");
  const files = fs.readdirSync("./log/v1/feed");
  for (const file of files) {
    const datas = await JSON.parse(
      fs.readFileSync(`./log/v1/feed/${file}.log`)
    ).data;
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
          .setColor(config.color.s)
          .setTimestamp();
        channel.send({ embeds: [embed] });
        datas[index].lastDate = items[0].pubDate;
        functions.logUtils.loging(datas, `v1/feed/${file}`);
      }
    });
  }
});

client.login(config.token);

// エラー処理 (これ入れないとエラーで落ちる。本当は良くないかもしれない)
process.on("uncaughtException", (error) => {
  console.error(
    `[${functions.timeUtils.timeToJST(Date.now(), true)}] ${error.stack}`
  );
  const embed = new EmbedBuilder()
    .setTitle("ERROR - uncaughtException")
    .setDescription("```\n" + error.stack + "\n```")
    .setColor(config.color.e)
    .setTimestamp();
  client.channels
    .fetch(config.logch.error)
    .then((c) => c.send({ embeds: [embed] }));
});

process.on("unhandledRejection", (reason, promise) => {
  console.error(
    `\u001b[31m[${functions.timeUtils.timeToJST(
      Date.now(),
      true
    )}] ${reason}\u001b[0m\n`,
    promise
  );
  const embed = new EmbedBuilder()
    .setTitle("ERROR - unhandledRejection")
    .setDescription("```\n" + reason + "\n```")
    .setColor(config.color.e)
    .setTimestamp();
  client.channels
    .fetch(config.logch.error)
    .then((c) => c.send({ embeds: [embed] }));
});
