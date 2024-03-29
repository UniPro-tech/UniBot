const Discord = require("discord.js");
module.exports = {
  name: "ready", // イベント名
  async execute(client) {
    const log = await client.func.readLog("v1/conf/status");
    const add = require(`../system/add.js`);

    add.addCmd(client.conf);
    console.log(`on:${log.onoff},play:${log.playing},status:${log.status}`);

    if (log.onoff == 'on') {

      let activityOpt = {};
      switch (log.type) {
        case "WATCHING":
          activityOpt.type = Discord.ActivityType.Watching;
          break;

        case "COMPETING":
          activityOpt.type = Discord.ActivityType.Competing;
          break

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

      client.user.setActivity(log.playing, activityOpt);
      if (log.status == 'Discord Android') {
        client.ws = { properties: { $browser: 'Discord Android' } };
      } else {
        //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
        client.user.setStatus(log.status);
      }
    } else {
      client.shard.fetchClientValues('guilds.cache.size')
        .then(result => {
          client.user.setActivity(`Servers: ${result}`);
        });
      //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
      client.user.setStatus('online');
    }
    client.channels.cache.get(client.conf.logch.ready).send("Discordログインしました!");
    console.log(`[${client.func.timeToJST(Date.now(), true)}] Logged in as ${client.user.tag}!`);
  }
}