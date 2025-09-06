import {
  Client,
  GatewayIntentBits,
  Collection,
  Partials,
  EmbedBuilder,
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
  ModalSubmitCollector,
  StringSelectMenuCollector,
} from "@/lib/collecter";
import { ChatInputCommand } from "./executors/types/ChatInputCommand";
import { StringSelectMenu } from "./executors/types/StringSelectMenu";
import { Button } from "./executors/types/Button";
import { LogContext, Logger, Transporter } from "@unipro-tech/node-logger";
import { AsyncLocalStorage } from "async_hooks";

export const loggingSystem = new Logger("unibot", [
  ...(process.env.NODE_ENV === "development"
    ? [Transporter.PinoPrettyTransporter()]
    : [Transporter.ConsoleTransporter()]),
]);

export const ALStorage = new AsyncLocalStorage<LogContext & Record<string, any>>();

const client = new Client({
  intents: [
    GatewayIntentBits.Guilds,
    GatewayIntentBits.GuildMessages,
    GatewayIntentBits.MessageContent,
    GatewayIntentBits.GuildVoiceStates,
  ],
  partials: [Partials.Channel],
});

import { Agenda } from "@hokify/agenda";
import { ModalSubmitCommand } from "./executors/types/ModalSubmit";
if (!process.env.DATABASE_URL) {
  throw new Error("DATABASE_URL is not defined in environment variables.");
}
export const agenda = new Agenda({ db: { address: process.env.DATABASE_URL } });

import jobManager from "./lib/jobManager";
import { nanoid } from "nanoid";

// Attach utilities and config to client
client.agenda = agenda;
client.config = config;
client.functions = { timeUtils, logUtils, jobManager };
client.fs = fs;

// Setup interaction executor collections
client.interactionExecutorsCollections = {
  chatInputCommands: new Collection<string, ChatInputCommand>(),
  stringSelectMenus: new Collection<string, StringSelectMenu>(),
  messageContextMenuCommands: new Collection<string, ChatInputCommand>(),
  buttons: new Collection<string, Button>(),
  modalSubmitCommands: new Collection<string, ModalSubmitCommand>(),
};

agenda.on("ready", async () => {
  const logger = loggingSystem.getLogger({ function: "agenda" });
  logger.info("Agenda started successfully.");
});

agenda.on("error", (error) => {
  const logger = loggingSystem.getLogger({ function: "agenda" });
  logger.error({ stack_trace: error.stack, error: error.message }, "Agenda an error occurred.");
});

agenda.define("purge agenda", async (job, done) => {
  const logger = loggingSystem.getLogger({ function: "agenda" });
  const jobs = await agenda.jobs();
  jobs.forEach((job) => {
    if (job.attrs.nextRunAt == null) logger.info("This job is finished and can be removed.");
    if (job.attrs.nextRunAt == null) job.remove();
  });
  done();
});

// Register collectors
ChatInputCommandCollector(client);
StringSelectMenuCollector(client);
MessageContextMenuCommandCollector(client);
ButtonCollector(client);
ModalSubmitCollector(client);

// Dynamically load event files
const eventDir = path.resolve(__dirname, "events");
const eventFiles = fs
  .readdirSync(eventDir)
  .filter((file) => (file.endsWith(".ts") && !file.endsWith(".d.ts")) || file.endsWith(".js"));

for (const file of eventFiles) {
  const logger = loggingSystem.getLogger({ function: "eventLoader" });
  const event = require(path.join(eventDir, file));
  const handler = (...args: any[]) => event.execute(...args, client);
  try {
    event.once
      ? client.once(event.name, async (...args: any[]) => {
          const ctx: LogContext = {
            trace_id: nanoid(),
            request_id: nanoid(),
          };
          ALStorage.run(ctx, () => {
            handler(...args);
          });
        })
      : client.on(event.name, (...args: any[]) => {
          const ctx: LogContext = {
            trace_id: nanoid(),
            request_id: nanoid(),
          };
          ALStorage.run(ctx, () => {
            handler(...args);
          });
        });
  } catch (error) {
    logger.error(
      { stack_trace: (error as Error).stack, extra_context: { event, file } },
      "Failed to load event."
    );
  }
}

// Error handling
const sendErrorEmbed = async (title: string, description: string) => {
  const logger = loggingSystem.getLogger({ function: "errorHandler" });
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
      logger.error(
        { extra_context: { channelId: config.logch.error, channel } },
        "Error log channel not found or is not a text channel."
      );
    }
  } catch (error) {
    logger.error({ stack_trace: (error as Error).stack, error }, (error as Error).message);
  }
};

process.on("uncaughtException", (error) => {
  const logger = loggingSystem.getLogger({ function: "errorHandler" });
  logger.error({ stack_trace: (error as Error).stack, error }, (error as Error).message);
  sendErrorEmbed("ERROR - uncaughtException", error.stack || String(error));
});

process.on("unhandledRejection", (reason: any, promise) => {
  const logger = loggingSystem.getLogger({ function: "errorHandler" });
  let reasonText = "";
  if (reason instanceof Error) {
    reasonText = reason.stack || reason.message;
  } else {
    reasonText = typeof reason === "object" ? JSON.stringify(reason, null, 2) : String(reason);
  }
  logger.error({ extra_context: { reason, promise } }, `An unhandledRejection occurred`);
  sendErrorEmbed("ERROR - unhandledRejection", reasonText);
});

client.login(config.token);
