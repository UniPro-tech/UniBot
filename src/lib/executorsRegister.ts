import { ChatInputCommand } from "@/executors/types/Command";
import { REST } from "@discordjs/rest";
import { RESTPostAPIChatInputApplicationCommandsJSONBody, Routes } from "discord-api-types/v10";
import { Client, SlashCommandBuilder } from "discord.js";
import fs from "fs";
import path from "path";

/**
 * @param {Client} client
 */
export const chatInputCommandsRegister = async (client: Client) => {
  console.log(`\u001b[32m===Pushing ChatInputCommand Data===\u001b[0m`);
  const config = client.config;
  const token = config.token;
  const rest = new REST({ version: "10" }).setToken(token);

  const testGuild = config.dev.testGuild;

  let command_int = 0;
  const globalCommands = [] as RESTPostAPIChatInputApplicationCommandsJSONBody[];
  const adminGuildCommands = [] as RESTPostAPIChatInputApplicationCommandsJSONBody[];
  const commandFolders = fs.readdirSync(path.resolve(__dirname, `../executors/chatInputCommands`));

  function cmdToArray(
    array: RESTPostAPIChatInputApplicationCommandsJSONBody[],
    command: ChatInputCommand,
    file: string,
    notice = ""
  ) {
    try {
      array.push((command.data as SlashCommandBuilder).toJSON());
      command_int++;
      console.log(`${notice} ${file} has been added.`);
    } catch (error) {
      console.error(`${notice} An Error Occured in ${file} \nエラー内容\n ${error}`);
    }
  }

  async function putToDiscord(
    array: RESTPostAPIChatInputApplicationCommandsJSONBody[],
    guild: undefined | string = undefined
  ) {
    if (guild) {
      await rest.put(Routes.applicationGuildCommands(client.application?.id as string, guild), {
        body: array,
      });
    } else {
      await rest.put(Routes.applicationCommands(client.application?.id as string), {
        body: array,
      });
    }
  }

  for (const folder of commandFolders) {
    console.log(`[Init]Adding ${folder} commands...`);
    const commandFiles = fs
      .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${folder}`))
      .filter((file) => file.endsWith(".js") || (file.endsWith(".ts") && !file.endsWith(".d.ts")));
    for (const file of commandFiles) {
      const command = require(path.resolve(
        __dirname,
        `../executors/chatInputCommands/${folder}/${file}`
      )) as ChatInputCommand;
      if (command.adminGuildOnly) {
        cmdToArray(adminGuildCommands, command, file, "[Admin Slash Command]");
        continue;
      }
      //if (command.onlyCommand) continue;
      cmdToArray(globalCommands, command, file, "[Global Slash Command]");
    }
    console.log(`[Init]${folder} added.`);
  }

  try {
    console.log(`[Init]Registering ${command_int}...`);

    //Admin
    putToDiscord(adminGuildCommands, testGuild);
    console.log(`[Init]Registered Admin Guild Slash Commands.`);

    //Global
    putToDiscord(globalCommands);
    console.log(`[Init]Registered Global Slash Commands.`);

    console.log(`[Init]Registered All Slash Commands.`);
  } catch (error) {
    console.error("[error]", error);
  }
};

/**
 * @param {Client} client
 */
export const messageContextMenuCommandsRegister = async (client: Client) => {
  console.log(`\u001b[32m===Pushing MessageContextMenuCommand Data===\u001b[0m`);
  const config = client.config;
  const token = config.token;
  const rest = new REST({ version: "10" }).setToken(token);

  const testGuild = config.dev.testGuild;

  let command_int = 0;
  const globalCommands = [] as RESTPostAPIChatInputApplicationCommandsJSONBody[];
  const adminGuildCommands = [] as RESTPostAPIChatInputApplicationCommandsJSONBody[];

  function cmdToArray(
    array: RESTPostAPIChatInputApplicationCommandsJSONBody[],
    command: ChatInputCommand,
    file: string,
    notice = ""
  ) {
    try {
      array.push((command.data as SlashCommandBuilder).toJSON());
      command_int++;
      console.log(`${notice} ${file} has been added.`);
    } catch (error) {
      console.error(`${notice} An Error Occured in ${file} \nエラー内容\n ${error}`);
    }
  }

  async function putToDiscord(
    array: RESTPostAPIChatInputApplicationCommandsJSONBody[],
    guild: undefined | string = undefined
  ) {
    if (guild) {
      await rest.put(Routes.applicationGuildCommands(client.application?.id as string, guild), {
        body: array,
      });
    } else {
      await rest.put(Routes.applicationCommands(client.application?.id as string), {
        body: array,
      });
    }
  }

  console.log(
    `[Init]Adding ${path.resolve(__dirname, `../executors/messageContextMenuCommands`)} commands...`
  );
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/messageContextMenuCommands`))
    .filter((file) => file.endsWith(".js") || (file.endsWith(".ts") && !file.endsWith(".d.ts")));
  for (const file of commandFiles) {
    const command = require(path.resolve(
      __dirname,
      `../executors/messageContextMenuCommands/${file}`
    )) as ChatInputCommand;
    if (command.adminGuildOnly) {
      cmdToArray(adminGuildCommands, command, file, "[Admin]");
      continue;
    }
    //if (command.onlyCommand) continue;
    cmdToArray(globalCommands, command, file, "[Global]");
  }

  try {
    console.log(`[Init]Registering ${command_int}...`);

    //Admin
    putToDiscord(adminGuildCommands, testGuild);
    console.log(`[Init]Registered Admin Guild Slash Commands.`);

    //Global
    putToDiscord(globalCommands);
    console.log(`[Init]Registered Global Slash Commands.`);

    console.log(`[Init]Registered All Slash Commands.`);
  } catch (error) {
    console.error("[error]", error);
  }
};
