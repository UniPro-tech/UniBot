import {
  Client,
  GatewayIntentBits,
  Collection,
  Partials,
  EmbedBuilder,
  Channel,
  TextChannel,
} from "discord.js";
import fs from "fs";
import path from "path";
import config from "@/config";
import timeUtils from "@/lib/timeUtils";
import logUtils from "@/lib/dataUtils";
import {
  ButtonCollector,
  ChatInputCommandCollector,
  MessageContextMenuCommandCollector,
  StringSelectMenuCollector,
} from "@/lib/collecter";
import { ChatInputCommand } from "./executors/types/ChatInputCommand";
import { StringSelectMenu } from "./executors/types/StringSelectMenu";
import { Button } from "./executors/types/Button";

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildVoiceStates,
  ],
  partials: [Partials.Channel],
});

// Attach utilities and config to client
client.config = config;
client.functions = { timeUtils, logUtils };
client.fs = fs;

// Setup interaction executor collections
client.interactionExecutorsCollections = {
  chatInputCommands: new Collection<string, ChatInputCommand>(),
  stringSelectMenus: new Collection<string, StringSelectMenu>(),
  messageContextMenuCommands: new Collection<string, ChatInputCommand>(),
  buttons: new Collection<string, Button>(),
};

// Register collectors
ChatInputCommandCollector(client);
StringSelectMenuCollector(client);
MessageContextMenuCommandCollector(client);
ButtonCollector(client);

// Dynamically load event files
const eventDir = path.resolve(__dirname, "events");
const eventFiles = fs
  .readdirSync(eventDir)
  .filter((file) => (file.endsWith(".ts") && !file.endsWith(".d.ts")) || file.endsWith(".js"));

for (const file of eventFiles) {
  const event = require(path.join(eventDir, file));
  const handler = (...args: any[]) => event.execute(...args, client);
  try {
    event.once ? client.once(event.name, handler) : client.on(event.name, handler);
  } catch (error) {
    console.error(
      `\u001b[31m[${client.functions.timeUtils.timeToJSTstamp(Date.now(), true)}]\u001b[0m\n`,
      error
    );
  }
}

// Error handling
const sendErrorEmbed = async (title: string, description: string) => {
  const embed = new EmbedBuilder()
    .setTitle(title)
    .setDescription("```\n" + description + "\n```")
    .setColor(config.color.error)
    .setTimestamp();
  try {
    const channel = await client.channels.fetch(config.logch.error);
    if (channel && channel instanceof TextChannel) {
      channel.send({ embeds: [embed] });
    } else {
      console.error("Error: Log Channel is invalid.");
    }
  } catch (err) {
    console.error("Failed to send error embed:", err);
  }
};

process.on("uncaughtException", (error) => {
  const timestamp = client.functions.timeUtils.timeToJSTstamp(Date.now(), true);
  console.error(`[${timestamp}] ${error.stack}`);
  sendErrorEmbed("ERROR - uncaughtException", error.stack || String(error));
});

process.on("unhandledRejection", (reason: any, promise) => {
  const timestamp = client.functions.timeUtils.timeToJSTstamp(Date.now(), true);
  let reasonText = "";
  if (reason instanceof Error) {
    reasonText = reason.stack || reason.message;
  } else {
    reasonText = typeof reason === "object" ? JSON.stringify(reason, null, 2) : String(reason);
  }
  console.error(`\u001b[31m[${timestamp}] ${reasonText}\u001b[0m`, promise);
  sendErrorEmbed("ERROR - unhandledRejection", reasonText);
});

client.login(config.token);
