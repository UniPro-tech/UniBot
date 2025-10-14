import {
  Client,
  EmbedBuilder,
  ActivityType,
  ActivityOptions,
  TextChannel,
  PresenceStatusData,
} from "discord.js";
import { registerAllCommands } from "@/lib/executorsRegister";
import path from "path";
import { redefineJobs } from "@/lib/jobManager";
import { ALStorage, loggingSystem } from "..";

export const name = "ready";

export const execute = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "ready" });
  const logFile = await client.functions.logUtils.readConfig("status");
  await registerAllCommands(client);

  logger.info({ extra_context: { log_file: logFile } }, "Bot is ready");

  if (!client.user) {
    logger.error({ extra_context: { service: "ready" } }, "Client user is not defined");
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

    // Use setPresence for clearer semantics and to avoid ambiguity with overloaded setActivity
    client.user.setPresence({
      activities: [
        { name: String(activityOpt.name), type: activityOpt.type, url: activityOpt.url },
      ],
      status: (logFile.status as PresenceStatusData) || "online",
    });
  } else {
    client.user.setPresence({
      activities: [{ name: `Servers: ${client.guilds.cache.size}`, type: ActivityType.Playing }],
      status: "online",
    });
  }

  // Fetch the channel to avoid relying on cache (channel may not be cached yet)
  const channel = await client.channels.fetch(client.config.logch.ready).catch(() => null);
  if (!channel || !(channel instanceof TextChannel)) {
    logger.error(
      { extra_context: { channel: client.config.logch.ready, fetched: !!channel } },
      "Log channel is not defined or not a text channel"
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

  await channel.send({ embeds: [embed] });

  await client.agenda.start();
  await client.agenda.now("purge agenda");
  client.agenda.every("0 0 * * *", "purge agenda");
  await redefineJobs(client);

  logger.info("Ready event processing completed");
};

export default {
  name,
  execute,
};
