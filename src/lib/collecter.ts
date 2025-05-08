import { Client, Collection } from "discord.js";
import path from "path";
import fs from "fs";

export const CommandCollector = async (client: Client) => {
  console.log(`\u001b[32m===Load Command Executing Data===\u001b[0m`);
  client.commands = new Collection();
  const commandFolders = fs.readdirSync(path.resolve(__dirname, `../commands`));
  for (const folder of commandFolders) {
    console.log(`\u001b[32m[Init]Loading ${folder} commands\u001b[0m`);
    const commandFiles = fs
      .readdirSync(path.resolve(__dirname, `../commands/${folder}`))
      .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
    for (const file of commandFiles) {
      console.debug(`dir:${folder},file:${file}`);
      const command = require(path.resolve(__dirname, `../commands/${folder}/${file}`));
      try {
        client.commands.set(command.data.name, command);
        console.log(`[Init]${command.data.name} has been loaded.`);
      } catch (error) {
        console.error(`[error]An Error Occured in ${command.data.name}\nDatails:\n ${error}`);
      }
    }
    console.log(`\u001b[32m${folder} has been loaded\u001b[0m`);
  }
};

export const StringSelectMenuCollector = async (client: Client) => {
  console.log(`\u001b[32m===Load Command Executing Data===\u001b[0m`);
  client.stringSelectMenus = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../selectMenus/string`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = require(path.resolve(__dirname, `../selectMenus/string/${file}`));
    try {
      client.stringSelectMenus.set(menuDefine.name, menuDefine);
      console.log(`[Init]${menuDefine.name} has been loaded.`);
    } catch (error) {
      console.error(`[error]An Error Occured in ${menuDefine.name}\nDetails:\n ${error}`);
    }
  }
};
