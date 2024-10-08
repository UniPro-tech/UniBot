const Discord = require("discord.js");
module.exports = {
  name: "ready", // イベント名
  /**
   * Event handler for the "ready" event.
   *
   * @param {Discord.Client} client - The Discord client object.
   * @returns {Promise<void>} - A promise that resolves when the execution is complete.
   */
  async execute(client) {
    const logFile = await client.func.logUtils.readLog("v1/conf/status");
    console.log(logFile);
    const commandUtils = require(`../lib/commandUtils.js`);
    commandUtils.addCmd(client);
    console.debug(
      `[debug] on:${logFile?.onoff},play:${logFile?.playing},status:${logFile?.status}`
    );
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
          activityOpt.url = statusDescription;
          activityOpt.name = "Youtube";
          break;

        case "CUSTOM":
          activityOpt.type = Discord.ActivityType.Custom;
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
    if (channel) {
      channel.send("Discordログインしました!");
    }
    console.log(
      `[${client.func.timeUtils.timeToJST(Date.now(), true)} info] Logged in as ${
        client.user.tag
      }!`
    );
  },
};
