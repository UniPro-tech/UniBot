import { Client, Collection } from "discord.js";
import path from "path";
import fs from "fs";
import { ALStorage, loggingSystem } from "..";

// TODO: Collectionの型を指定する

export const ChatInputCommandCollector = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "ChatInputCommandCollector" });
  logger.info("Load ChatInputCommand Executing Data");
  client.interactionExecutorsCollections.chatInputCommands = new Collection();
  const commandFolders = fs.readdirSync(path.resolve(__dirname, `../executors/chatInputCommands`));
  for (const folder of commandFolders) {
    logger.info({ extra_context: { folder } }, `Started loading commands`);
    const commandFiles = fs
      .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${folder}`))
      .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
    for (const file of commandFiles) {
      const command = require(path.resolve(
        __dirname,
        `../executors/chatInputCommands/${folder}/${file}`
      ));
      try {
        client.interactionExecutorsCollections.chatInputCommands.set(command.data.name, command);
        logger.info(
          { extra_context: { commandName: command.data.name } },
          `Command has been loaded.`
        );
      } catch (error) {
        logger.error(
          {
            extra_context: { commandName: command.data.name },
            error: error,
            stack_trace: (error as Error).stack,
          },
          (error as Error).message
        );
      }
    }
    logger.info({ service: "Init" }, `Loaded ${folder} commands`);
  }
  logger.info({ service: "Init" }, "ChatInputCommand Executing Data Loaded");
};

export const StringSelectMenuCollector = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "StringSelectMenuCollector" });
  logger.info("Load StringSelectMenu Executing Data");
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
      logger.info({ extra_context: { command: menuDefine.name } }, `Command has been loaded.`);
    } catch (error) {
      logger.error(
        { stack_trace: (error as Error).stack, error: (error as Error).message },
        (error as Error).message
      );
    }
  }
  logger.info("StringSelectMenu Executing Data Loaded");
};

export const MessageContextMenuCommandCollector = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({
    ...ctx,
    function: "MessageContextMenuCommandCollector",
  });
  logger.info("Load MessageContextMenuCommand Executing Data");
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
      logger.info({ extra_context: { command: menuDefine.name } }, `Command has been loaded.`);
    } catch (error) {
      logger.error(
        { stack_trace: (error as Error).stack, error, extra_context: { command: menuDefine.name } },
        (error as Error).message
      );
    }
  }
  logger.info("MessageContextMenuCommand Executing Data Loaded");
};

export const ButtonCollector = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "ButtonCollector" });
  logger.info("Load Button Executing Data");
  client.interactionExecutorsCollections.buttons = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/buttons`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(path.resolve(__dirname, `../executors/buttons/${file}`));
    try {
      client.interactionExecutorsCollections.buttons.set(menuDefine.name, menuDefine);
      logger.info({ extra_context: { command: menuDefine.name } }, `Command has been loaded.`);
    } catch (error) {
      logger.error(
        {
          extra_context: { command: menuDefine.name },
          stack_trace: (error as Error).stack,
          error: error as Error,
        },
        (error as Error).message
      );
    }
  }
  logger.info("Button Executing Data Loaded");
};

export const ModalSubmitCollector = async (client: Client) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "ModalSubmitCollector" });
  logger.info("Load Modals Executing Data");
  client.interactionExecutorsCollections.modalSubmitCommands = new Collection();
  const commandFiles = fs
    .readdirSync(path.resolve(__dirname, `../executors/modals`))
    .filter((file) => file.endsWith(".ts") && !file.endsWith(".d.ts"));
  for (const file of commandFiles) {
    const menuDefine = await import(path.resolve(__dirname, `../executors/modals/${file}`));
    try {
      client.interactionExecutorsCollections.modalSubmitCommands.set(menuDefine.name, menuDefine);
      logger.info({ extra_context: { command: menuDefine.name } }, `Command has been loaded.`);
    } catch (error) {
      logger.error(
        {
          stack_trace: (error as Error).stack,
          error: error as Error,
          extra_context: { command: menuDefine.name },
        },
        (error as Error).message
      );
    }
  }
  logger.info("Modal Executing Data Loaded");
};
