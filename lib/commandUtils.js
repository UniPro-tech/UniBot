const { REST } = require("@discordjs/rest");
const { Routes } = require("discord-api-types/v9");
const { Client, Collection } = require("discord.js");
const fs = require("fs");

module.exports = {
  /**
   * @param {Client} client
   */
  async addCmd(client) {
    console.log(`\u001b[32m===Pushing Command Data===\u001b[0m`);
    const config = client.conf;
    const token = config.token;
    const rest = new REST({ version: "9" }).setToken(token);

    const testGuild = config.dev.testGuild;

    let command_int = 0;
    const globalCommands = [];
    const adminGuildCommands = [];
    const commandFolders = fs.readdirSync(`${process.cwd()}/commands`);

    function cmdToArray(array, command, file, notice = "") {
      try {
        array.push(command.data.toJSON());
        command_int++;
        console.log(`${notice} ${file} has been added.`);
      } catch (error) {
        console.error(
          `${notice} An Error Occured in ${file} \nエラー内容\n ${error}`
        );
      }
    }

    async function putToDiscord(array, guild = false) {
      if (guild) {
        await rest.put(
          Routes.applicationGuildCommands(client.application.id, guild),
          {
            body: array,
          }
        );
      } else {
        await rest.put(Routes.applicationCommands(client.application.id), {
          body: array,
        });
      }
    }

    for (const folder of commandFolders) {
      console.log(`[Init]Adding ${folder} commands...`);
      const commandFiles = fs
        .readdirSync(`${process.cwd()}/commands/${folder}`)
        .filter((file) => file.endsWith(".js"));
      for (const file of commandFiles) {
        const command = require(`${process.cwd()}/commands/${folder}/${file}`);
        if (command.adminGuildOnly) {
          cmdToArray(adminGuildCommands, command, file, "[Admin]");
          continue;
        }
        if (command.onlyCommand) continue;
        cmdToArray(globalCommands, command, file, "[Global]");
      }
      console.log(`[Init]${folder} added.`);
    }

    (async () => {
      try {
        console.log(
          `[Init]Registering ${command_int}...`
        );

        //Admin
        putToDiscord(adminGuildCommands, testGuild);
        console.log(
          `[Init]Registered Admin Guild Slash Commands.`
        );

        //Global
        putToDiscord(globalCommands);
        console.log(
          `[Init]Registered Global Slash Commands.`
        );

        console.log(
          `[Init]Registered All Slash Commands.`
        );
      } catch (error) {
        console.error(error);
      }
    })();
  },
  /**
   * Handles the command handling process.
   *
   * @param {Object} client - The client object.
   * @param {Object} fs - The fs object.
   * @param {Object} Collection - The Collection object.
   * @returns {Promise<void>} - A promise that resolves when the command handling is complete.
   */
  async handling(client) {
    console.log(`\u001b[32m===Load Command Executing Data===\u001b[0m`);
    client.commands = new Collection();
    const commandFolders = fs.readdirSync(`${process.cwd()}/commands`);
    for (const folder of commandFolders) {
      console.log(`\u001b[32m[Init]Loading ${folder} commands\u001b[0m`);
      const commandFiles = fs
        .readdirSync(`${process.cwd()}/commands/${folder}`)
        .filter((file) => file.endsWith(".js"));
      for (const file of commandFiles) {
        console.debug(`dir:${folder},file:${file}`);
        const command = require(`${process.cwd()}/commands/${folder}/${file}`);
        try {
          client.commands.set(command.data.name, command);
          console.log(
            `[Init]${
              command.data.name
            } has been loaded.`
          );
        } catch (error) {
          console.error(
            `[error]An Error Occured in ${
              command.data.name
            }\nDatails:\n ${error}`
          );
        }
      }
      console.log(`\u001b[32m${folder} has been loaded\u001b[0m`);
    }
  },
  /**
   * Adds subcommands to the provided data object.
   *
   * @param {string} name - The name of the command.
   * @param {object} data - The data object to add subcommands to.
   * @returns {object} - The modified data object with added subcommands.
   */
  addSubCommand(name, data) {
    console.log(`\u001b[32m[Init]Adding ${name}'s SubCommands\u001b[0m`);
    const commandFiles = fs
      .readdirSync(`${process.cwd()}/commands/${name}`)
      .filter((file) => file.endsWith(".js"));
    for (const file of commandFiles) {
      const command = require(`${process.cwd()}/commands/${name}/${file}`);
      if (command.subCommandGroup) {
        data.addSubcommandGroup(command.data);
      } else data.addSubcommand(command.data);

      console.log(`[Subcommand]${file} has been added.`);
    }
    console.log(`\u001b[32m[Init]Added ${name}'s SubCommands\u001b[0m`);
    return data;
  },
  /**
   * Handles sub-commands for a given collection and name.
   *
   * @param {Collection} collection - The collection to store the sub-commands.
   * @param {string} name - The name of the sub-commands.
   * @returns {Promise<void>} - A promise that resolves when the sub-commands are handled.
   */
  subCommandHandling(name) {
    const collection = new Collection();
    console.log(`\u001b[32m===Load ${name}'s SubCommand Executing Data===\u001b[0m`);
    const commandFiles = fs
      .readdirSync(`${process.cwd()}/commands/${name}`)
      .filter((file) => file.endsWith(".js"));
    for (const file of commandFiles) {
      console.log(`dir:${name},file:${file}`);
      const command = require(`${process.cwd()}/commands/${name}/${file}`);
      try {
        collection.set(command.data.name, command);
        console.log(
          `[Subcommand]${
            command.data.name
          } has been loaded.`
        );
      } catch (error) {
        console.log(
          `\u001b[31m[error]An Error Occured in ${
            command.data.name
          }\nDatails:\n ${error}\u001b[0m`
        );
      }
    }
    console.log(`\u001b[32m===${name} loaded===\u001b[0m`);
    return collection;
  },
};
