import { ChatInputCommand } from "@/executors/types/ChatInputCommand";
import { Collection, SlashCommandBuilder, SlashCommandSubcommandBuilder } from "discord.js";
import fs from "fs";
import path from "path";

/**
 * Adds subcommands to the provided data object.
 *
 * @param {string} name - The name of the command.
 * @param {object} data - The data object to add subcommands to.
 * @returns {object} - The modified data object with added subcommands.
 */
export const addSubCommand = (name: string, data: SlashCommandBuilder) => {
  console.log(`\u001b[32m[Init]Adding ${name}'s SubCommands\u001b[0m`);
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

    console.log(`[Subcommand]${file} has been added.`);
  }
  console.log(`\u001b[32m[Init]Added ${name}'s SubCommands\u001b[0m`);
  return data;
};

/**
 * Handles sub-commands for a given collection and name.
 *
 * @param {Collection} collection - The collection to store the sub-commands.
 * @param {string} name - The name of the sub-commands.
 * @returns {Promise<void>} - A promise that resolves when the sub-commands are handled.
 */
export const subCommandHandling = (name: string) => {
  const collection = new Collection<string, ChatInputCommand>();
  console.log(`\u001b[32m===Load ${name}'s SubCommand Executing Data===\u001b[0m`);
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
      console.log(`[Subcommand]${command.data.name} has been loaded.`);
    } catch (error) {
      console.log(
        `\u001b[31m[error]An Error Occured in ${command.data.name}\nDatails:\n ${error}\u001b[0m`
      );
    }
  }
  console.log(`\u001b[32m===${name} loaded===\u001b[0m`);
  return collection;
};

export default {
  addSubCommand,
  subCommandHandling,
};
