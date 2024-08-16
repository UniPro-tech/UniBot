const Discord = require("discord.js");
module.exports = {
  name: "ready", // イベント名
  async execute(client) {
    const logFile = await client.func.readLog("v1/conf/status");
    console.log(logFile);
    const add = require(`../system/commandRegister`);

    add.addCmd(client.conf);
    console.log(`on:${logFile?.onoff},play:${logFile?.playing},status:${logFile?.status}`);
    if (logFile?.onoff == "on") {
      let activityOpt = {};
      switch (logFile.type) {
        case "WATCHING":
          activityOpt.type = Discord.ActivityType.Watching;
          break;

        case "COMPETING":
          activityOpt.type = Discord.ActivityType.Competing;
          break;

        case "LISTENING":
          activityOpt.type = Discord.ActivityType.Listening;
          break;

        case "STREAMING":
          activityOpt.type = Discord.ActivityType.Streaming;
          break;

        default:
          activityOpt.type = Discord.ActivityType.Playing;
          break;
      }

      client.user.setActivity(logFile.playing, activityOpt);
      if (logFile.status == "Discord Android") {
        client.ws = { properties: { $browser: "Discord Android" } };
      } else {
        //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
        client.user.setStatus(logFile.status);
      }
    } else {
      client.user.setActivity(`Servers: ${client.guilds.cache.size}`);
      //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
      client.user.setStatus("online");
    }
    const channel = client.channels.cache.get(client.conf.logch.ready);
    channel.send("Discordログインしました!");
    console.log(
      `[${client.func.timeToJST(Date.now(), true)}] Logged in as ${client.user.tag}!`
    );
  },
};
