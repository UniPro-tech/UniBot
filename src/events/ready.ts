import Discord, { Client } from "discord.js";
import commandUtils from "@/lib/commandUtils";

export const name = "ready";
export const execute = async (client: Client) => {
  const logFile = await client.function.logUtils.read("v1/conf/status");
  commandUtils.addCommand(client);
  console.debug(
    `[debug] on:${logFile?.onoff},play:${logFile?.playing},status:${logFile?.status}`
  );
  if (!client.user) {
    console.error(
      `[error] [${client.function.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error] Client.user is undefined`
    );
    process.exit(1);
  }
  if (logFile?.onoff == "on") {
    let activityOpt: Discord.ActivityOptions = {
      name: logFile.playing,
      type: Discord.ActivityType.Playing,
      url: "",
    };
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
        //activityOpt.url = statusDescription;
        activityOpt.name = "Youtube";
        break;

      case "CUSTOM":
        activityOpt.type = Discord.ActivityType.Custom;
        break;

      default:
        activityOpt.type = Discord.ActivityType.Playing;
        break;
    }

    client.user?.setActivity(logFile.playing, activityOpt);
    /*if (logFile.status == "Discord Android") {
      client.ws = { properties: { $browser: "Discord Android" } };
    } else {
      //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
      client.user.setStatus(logFile.status);
    }*/
  } else {
    client.user?.setActivity(`Servers: ${client.guilds.cache.size}`);
    //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
    client.user?.setStatus("online");
  }
  const channel = client.channels.cache.get(client.config.logch.ready);
  if (!channel || !(channel instanceof Discord.TextChannel)) {
    console.error(
      `[error] [${client.function.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error] Log Channel is invalid`
    );
    return;
  } else {
    channel.send("Discordログインしました!");
  }
  console.log(
    `[${client.function.timeUtils.timeToJSTstamp(
      Date.now(),
      true
    )} info] Logged in as ${client.user.tag}!`
  );
};

export default {
  name,
  execute,
};
