import { ChatInputCommand } from "@/executors/types/ChatInputCommand";
import { StringSelectMenu } from "@/executors/types/StringSelectMenu";
import {
  Collection,
  SlashCommandBuilder,
  SlashCommandSubcommandBuilder,
  SlashCommandSubcommandGroupBuilder,
} from "discord.js";
import fs from "fs";
import path from "path";
import { loggingSystem } from "..";

/**
 * Adds subcommands to the provided data object.
 *
 * @param {string} name - The name of the command.
 * @param {object} data - The data object to add subcommands to.
 * @returns {object} - The modified data object with added subcommands.
 */
export const addSubCommand = (
  name: string,
  data: SlashCommandBuilder | SlashCommandSubcommandGroupBuilder
) => {
  const logger = loggingSystem.getLogger({ function: "addSubCommand" });
  logger.info({ extra_context: { command: name } }, `Adding SubCommands`);
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${name}`))
    .filter((file) => file.endsWith(".js") || file.endsWith(".ts"));
  for (const file of commandFiles) {
    const command = require(path.resolve(
      __dirname,
      `../executors/chatInputCommands/${name}/${file}`
    )) as ChatInputCommand;
    /*if (command.subCommandGroup) {
      data.addSubcommandGroup(command.data);
    } else */
    data.addSubcommand(command.data as SlashCommandSubcommandBuilder);

    logger.info({ extra_context: { command: command.data.name } }, "Subcommand has been added.");
  }
  logger.info({ extra_context: { command: name } }, `Added SubCommands`);
  return data;
};

/**
 * Handles sub-commands for a given collection and name.
 *
 * @param {Collection} collection - The collection to store the sub-commands.
 * @param {string} name - The name of the sub-commands.
 * @returns {Promise<void>} - A promise that resolves when the sub-commands are handled.
 */
export const subCommandHandling = (
  name: string,
  collection?: Collection<string, ChatInputCommand>
) => {
  const logger = loggingSystem.getLogger({ function: "subCommandHandling" });
  if (!collection) {
    collection = new Collection<string, ChatInputCommand>();
  }
  logger.info({ extra_context: { command: name } }, `Load SubCommand Executing Data`);
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${name}`))
    .filter((file) => file.endsWith(".js") || file.endsWith(".ts"));
  for (const file of commandFiles) {
    const command = require(path.resolve(
      __dirname,
      `../executors/chatInputCommands/${name}/${file}`
    )) as ChatInputCommand;
    try {
      collection.set(command.data.name, command);
      logger.info({ extra_context: { command: command.data.name } }, "Subcommand has been loaded.");
    } catch (error) {
      logger.error(
        {
          extra_context: { command: command.data.name },
          error,
          stack_trace: error instanceof Error ? error.stack : undefined,
        },
        error instanceof Error ? error.message : `An Error Occurred in ${command.data.name}`
      );
    }
  }
  logger.info({ service: "CommandUtils", name }, `${name} loaded`);
  return collection;
};

export const addSubCommandGroup = (name: string, data: SlashCommandBuilder) => {
  const logger = loggingSystem.getLogger({ function: "addSubCommandGroup" });
  logger.info({ extra_context: { command: name } }, `Adding SubCommandGroups`);
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${name}`))
    .filter((file) => (file.endsWith(".js") || file.endsWith(".ts")) && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const command = require(path.resolve(
      __dirname,
      `../executors/chatInputCommands/${name}/${file}`
    )) as ChatInputCommand;
    /*if (command.subCommandGroup) {
      data.addSubcommandGroup(command.data);
    } else */
    data.addSubcommandGroup(command.data as SlashCommandSubcommandGroupBuilder);

    logger.info(
      { extra_context: { command: command.data.name } },
      "SubCommandGroup has been added."
    );
  }
  logger.info({ extra_context: { command: name } }, `Added ${name}'s SubCommandGroups`);
  return data;
};

export const subSelectMenusHandling = (
  name: string,
  collection?: Collection<string, StringSelectMenu>
) => {
  if (!collection) {
    collection = new Collection<string, StringSelectMenu>();
  }
  const logger = loggingSystem.getLogger({ function: "subSelectMenusHandling" });
  logger.info({ extra_context: { command: name } }, `Load SelectMenu Executing Data`);
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/selectMenus/${name}`))
    .filter((file) => file.endsWith(".js") || file.endsWith(".ts"));
  for (const file of commandFiles) {
    const command = require(path.resolve(
      __dirname,
      `../executors/selectMenus/${name}/${file}`
    )) as StringSelectMenu;
    try {
      collection.set(command.name, command);
      logger.info({ extra_context: { command: command.name } }, "SelectMenu has been loaded.");
    } catch (error) {
      logger.error(
        {
          extra_context: { command: command.name },
          error,
          stack_trace: error instanceof Error ? error.stack : undefined,
        },
        error instanceof Error ? error.message : `An Error Occurred in ${command.name}`
      );
    }
  }
  logger.info({ extra_context: { command: name } }, `Command loaded`);
  return collection;
};

export default {
  addSubCommand,
  subCommandHandling,
  addSubCommandGroup,
};
