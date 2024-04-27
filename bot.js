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
      $device: "Replit Container",
    },
  },
});
const { exec } = require("child_process");
const cron = require("node-cron");
cron.schedule("0 */5 * * * *", () => {
  console.log("Ratecheck");
  exec("curl -v https://discord.com/api/v10", (err, stdout, stderr) => {
    if (err) {
      throw err;
    } else {
      if (stdout.match(/1015/)) {
        exec("sudo kill 1");
      }
    }
  });
});

const fs = require("fs");

const config = require("./config.js");
const functions = require("./function.js");

client.conf = config;
client.func = functions;
client.fs = fs;

const cmdH = require(`./system/command.js`);
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

// コマンドが来た時
client.on("interactionCreate", async (i) => {
  console.log(i.commandName);
  if (!i.isCommand()) return;
  const command = client.commands.get(i.commandName);
  if (!command) return;

  // DM専用コマンド
  if (command.guildOnly && !i.inGuild()) {
    const embed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このコマンドはDMでは実行できません。")
      .setColor(config.color.e);
    i.reply({ embeds: [embed] });
    return;
  }

  // 実行
  try {
    const logMsg = await command.execute(i, client);
    console.log(`[Run : ${i.commandName}]${logMsg}`);
    const log = new EmbedBuilder()
      .setTitle("コマンド実行ログ")
      .setDescription(`${i.user.tag}(${i.user.id}) がコマンドを実行しました。`)
      .setColor(config.color.s)
      .setTimestamp()
      .setThumbnail(i.user.displayAvatarURL({ dynamic: true }))
      .addFields([
        { name: "コマンド", value: "```\n" + i.toString() + "\n```" },
        {
          name: "実行サーバー",
          value:
            "```\n" +
            `${i.guild?.id ?? "DM"}(${i.guild?.id ?? "DM"})` +
            "\n```",
        },
        {
          name: "実行ユーザー",
          value: "```\n" + `${i.user.tag}(${i.user.id})` + "\n```",
        },
      ])
      .setFooter({ text: String(i.id) });
    client.channels
      .fetch(config.logch.command)
      .then((c) => c.send({ embeds: [log] }));
  } catch (error) {
    console.error(error);
    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.e)
      .setTimestamp();
    client.channels.fetch(config.logch.error).then({ embeds: [logEmbed] });
    const iEmbed = new EmbedBuilder()
      .setTitle("すみません、エラーが発生しました...")
      .setDescription("```\n" + error + "\n```")
      .setColor(config.color.e)
      .setTimestamp();
    i.reply(config.logch.error).then((c) =>
      config.logch.command.send({ embeds: [iEmbed] })
    );
  }
});
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
