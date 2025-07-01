import { Client, EmbedBuilder, ActivityType, ActivityOptions, TextChannel } from "discord.js";
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

  if (logFile?.onoff === "on") {
    const activityOpt: ActivityOptions = {
      name: logFile.playing,
      type: ActivityType.Playing,
      url: "",
    };

    switch (logFile.type) {
      case "WATCHING":
        activityOpt.type = ActivityType.Watching;
        break;
      case "COMPETING":
        activityOpt.type = ActivityType.Competing;
        break;
      case "LISTENING":
        activityOpt.type = ActivityType.Listening;
        break;
      case "STREAMING":
        activityOpt.type = ActivityType.Streaming;
        activityOpt.name = "Youtube";
        break;
      case "CUSTOM":
        activityOpt.type = ActivityType.Custom;
        break;
      default:
        activityOpt.type = ActivityType.Playing;
    }

    client.user.setActivity(logFile.playing, activityOpt);
  } else {
    client.user.setActivity(`Servers: ${client.guilds.cache.size}`);
    client.user.setStatus("online");
  }

  const channel = client.channels.cache.get(client.config.logch.ready);
  if (!channel || !(channel instanceof TextChannel)) {
    console.error(
      `[error] [${client.functions.timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error] Log Channel is invalid`
    );
    return;
  }

  const packageData = await import(path.resolve(__dirname, "../../package.json"));
  const embed = new EmbedBuilder()
    .setTitle("Bot Ready")
    .setColor(client.config.color.success)
    .addFields(
      { name: "Bot Name", value: client.user.tag, inline: true },
      { name: "Bot ID", value: client.user.id, inline: true },
      { name: "Guilds", value: `${client.guilds.cache.size} servers`, inline: true },
      { name: "Users", value: `${client.users.cache.size} users`, inline: true },
      { name: "Channels", value: `${client.channels.cache.size} channels`, inline: true }
    )
    .setFooter({
      text: `Bot Name: ${client.user.tag}`,
      iconURL: client.user.displayAvatarURL(),
    })
    .setTimestamp()
    .setThumbnail(client.user.displayAvatarURL({ size: 1024 }));

  if (packageData.description) embed.setDescription(packageData.description);
  if (packageData.version)
    embed.addFields({ name: "Version", value: packageData.version, inline: true });
  if (packageData.license)
    embed.addFields({ name: "License", value: packageData.license, inline: true });
  if (packageData.author)
    embed.addFields({ name: "Authors", value: packageData.author.name, inline: false });
  if (packageData.repository)
    embed.addFields({ name: "Repository", value: packageData.repository.url, inline: false });
  if (packageData.homepage) embed.setURL(packageData.homepage);

  if (packageData.contributors) {
    const contributors = packageData.contributors.map(
      (c: any) => `- [${c.name}](${c.url}) <[${c.email}](mailto:${c.email})>`
    );
    embed.addFields({ name: "Contributors", value: contributors.join("\n") });
  }

  console.debug(`[debug] Bot is ready and logged in as ${client.user.tag}`);
  await channel.send({ embeds: [embed] });

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
