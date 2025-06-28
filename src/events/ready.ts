import Discord, { Client } from "discord.js";
import { registerAllCommands } from "@/lib/executorsRegister";
import path from "path";

export const name = "ready";
export const execute = async (client: Client) => {
  const logFile = await client.functions.logUtils.readConfig("status");
  await registerAllCommands(client);
  console.debug(`[debug] on:${logFile?.onoff},play:${logFile?.playing},status:${logFile?.status}`);
  if (!client.user) {
    console.error(
      `[error] [${client.functions.timeUtils.timeToJSTstamp(
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
      `[error] [${client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error] Log Channel is invalid`
    );
    return;
  } else {
    const packageData = await import(path.resolve(__dirname, "../../../../package.json"));
    const embed = new Discord.EmbedBuilder()
      .setTitle("Bot Ready")
      .setDescription("Bot is ready and logged in successfully!")
      .setColor(client.config.color.success)
      .addFields([
        {
          name: "Bot Name",
          value: client.user.tag,
          inline: true,
        },
        {
          name: "Bot ID",
          value: client.user.id,
          inline: true,
        },
        {
          name: "Guilds",
          value: `${client.guilds.cache.size} servers`,
          inline: true,
        },
        {
          name: "Users",
          value: `${client.users.cache.size} users`,
          inline: true,
        },
        {
          name: "Channels",
          value: `${client.channels.cache.size} channels`,
          inline: true,
        },
      ])
      .setFooter({
        text: `Bot Name: ${client.user.tag}`,
        iconURL: client.user.displayAvatarURL(),
      })
      .setTimestamp();
    if (packageData.description) {
      embed.setDescription(packageData.description);
    }
    if (packageData.version) {
      embed.addFields({
        name: "Version",
        value: packageData.version,
        inline: true,
      });
    }
    if (packageData.license) {
      embed.addFields({
        name: "License",
        value: packageData.license,
        inline: true,
      });
    }
    if (packageData.author) {
      embed.addFields({
        name: "Authors",
        value: packageData.author.name,
        inline: false,
      });
    }
    if (packageData.repository) {
      embed.addFields({
        name: "Repository",
        value: packageData.repository.url,
        inline: false,
      });
    }
    if (packageData.homepage) {
      embed.setURL(packageData.homepage);
    }
    if (packageData.contributors) {
      let temp = [];
      for (let i = 0; i < packageData.contributors.length; i++) {
        temp[
          i
        ] = `- [${packageData.contributors[i].name}](${packageData.contributors[i].url}) <[${packageData.contributors[i].email}](mailto:${packageData.contributors[i].email})>`;
      }
      embed.addFields({ name: "Contributors", value: temp.join("\n") });
    }
    embed.setThumbnail(
      client.user?.displayAvatarURL({
        size: 1024,
      })
    );
    channel.send({ embeds: [embed] });
  }
  console.log(
    `[${client.functions.timeUtils.timeToJSTstamp(Date.now(), true)} info] Logged in as ${
      client.user.tag
    }!`
  );
};

export default {
  name,
  execute,
};
