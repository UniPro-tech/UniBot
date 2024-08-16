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
      $os: "Untitled OS",
      $browser: "Untitled Browser",
      $device: "K8s on Proxmox VE 6.4 in Sibainu",
    },
  },
});

const fs = require("fs");

const config = require("./config.js");
const functions = require("./function.js");

client.conf = config;
client.func = functions;
client.fs = fs;

const cmdH = require(`./system/command.js`);
client.commands = new Collection(); // Add this line to define the client.commands object
cmdH.handling(client, fs, Collection, config);
// イベントハンドリング
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
        `\u001b[31m[${yuika.timeToJST(Date.now(), true)}]\u001b[0m\n`,
        error
      );
    }
  } else {
    try {
      client.on(event.name, (...args) => event.execute(...args, client));
    } catch (error) {
      console.error(
        `\u001b[31m[${yuika.timeToJST(Date.now(), true)}]\u001b[0m\n`,
        error
      );
    }
  }
}

client.login(config.token);

// エラー処理 (これ入れないとエラーで落ちる。本当は良くないかもしれない)
process.on("uncaughtException", (error) => {
  console.error(`[${functions.timeToJST(Date.now(), true)}] ${error.stack}`);
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
    `\u001b[31m[${functions.timeToJST(Date.now(), true)}] ${reason}\u001b[0m\n`,
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
