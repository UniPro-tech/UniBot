import { Client, Collection } from "discord.js";
import path from "path";
import fs from "fs";

// TODO: Collectionの型を指定する

export const ChatInputCommandCollector = async (client: Client) => {
  console.info(`\u001b[32m===Load ChatInputCommand Executing Data===\u001b[0m`);
  client.interactionExecutorsCollections.chatInputCommands = new Collection();
  const commandFolders = fs.readdirSync(path.resolve(__dirname, `../executors/chatInputCommands`));
  for (const folder of commandFolders) {
    console.info(`\u001b[32m[Init]Loading ${folder} commands\u001b[0m`);
    const commandFiles = fs
      .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${folder}`))
      .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
    for (const file of commandFiles) {
      console.debug(`dir:${folder},file:${file}`);
      const command = require(path.resolve(
        __dirname,
        `../executors/chatInputCommands/${folder}/${file}`
      ));
      try {
        client.interactionExecutorsCollections.chatInputCommands.set(command.data.name, command);
        console.info(`[Init]${command.data.name} has been loaded.`);
      } catch (error) {
        console.error(`[error]An Error Occurred in ${command.data.name}\nDetails:\n ${error}`);
      }
    }
    console.info(`\u001b[32m${folder} has been loaded\u001b[0m`);
  }
  console.info(`\u001b[32m===ChatInputCommand Executing Data Loaded===\u001b[0m`);
};

export const StringSelectMenuCollector = async (client: Client) => {
  console.info(`\u001b[32m===Load StringSelectMenu Executing Data===\u001b[0m`);
  client.interactionExecutorsCollections.stringSelectMenus = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/selectMenus/string`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(
      path.resolve(__dirname, `../executors/selectMenus/string/${file}`)
    );
    try {
      client.interactionExecutorsCollections.stringSelectMenus.set(menuDefine.name, menuDefine);
      console.info(`[Init]${menuDefine.name} has been loaded.`);
    } catch (error) {
      console.error(`[error]An Error Occurred in ${menuDefine.name}\nDetails:\n ${error}`);
    }
  }
  console.info(`\u001b[32m===StringSelectMenu Executing Data Loaded===\u001b[0m`);
};

export const MessageContextMenuCommandCollector = async (client: Client) => {
  console.info(`\u001b[32m===Load MessageContextMenuCommand Executing Data===\u001b[0m`);
  client.interactionExecutorsCollections.messageContextMenuCommands = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/messageContextMenuCommands`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(
      path.resolve(__dirname, `../executors/messageContextMenuCommands/${file}`)
    );
    try {
      client.interactionExecutorsCollections.messageContextMenuCommands.set(
        menuDefine.name,
        menuDefine
      );
      console.info(`[Init]${menuDefine.name} has been loaded.`);
    } catch (error) {
      console.error(`[error]An Error Occurred in ${menuDefine.name}\nDetails:\n ${error}`);
    }
  }
  console.info(`\u001b[32m===MessageContextMenuCommand Executing Data Loaded===\u001b[0m`);
};

export const ButtonCollector = async (client: Client) => {
  console.info(`\u001b[32m===Load Button Executing Data===\u001b[0m`);
  client.interactionExecutorsCollections.buttons = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/buttons`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(path.resolve(__dirname, `../executors/buttons/${file}`));
    try {
      client.interactionExecutorsCollections.buttons.set(menuDefine.name, menuDefine);
      console.info(`[Init]${menuDefine.name} has been loaded.`);
    } catch (error) {
      console.error(`[error]An Error Occurred in ${menuDefine.name}\nDetails:\n ${error}`);
    }
  }
  console.info(`\u001b[32m===Button Executing Data Loaded===\u001b[0m`);
};

export const ModalSubmitCollector = async (client: Client) => {
  console.info(`\u001b[32m===Load Modals Executing Data===\u001b[0m`);
  client.interactionExecutorsCollections.modalSubmitCommands = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/modals`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(path.resolve(__dirname, `../executors/modals/${file}`));
    try {
      client.interactionExecutorsCollections.modalSubmitCommands.set(menuDefine.name, menuDefine);
      console.info(`[Init]${menuDefine.name} has been loaded.`);
    } catch (error) {
      console.error(`[error]An Error Occurred in ${menuDefine.name}\nDetails:\n ${error}`);
    }
  }
  console.info(`\u001b[32m===Modal Executing Data Loaded===\u001b[0m`);
};
