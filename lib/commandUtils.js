const { REST } = require("@discordjs/rest");
const { Routes } = require("discord-api-types/v9");
const { Client } = require("discord.js");
const fs = require("fs");

module.exports = {
  /**
   * @param {Client} client
   */
  async addCmd(client) {
    const config = client.config;
    const token = config.token;
    const rest = new REST({ version: "9" }).setToken(token);

    const testGuild = config.dev.testGuild;

    let command_int = 0;
    const globalCommands = [];
    const adminGuildCommands = [];
    const commandFolders = fs.readdirSync(`./commands`);

    function cmdToArray(array, command, file, notice = "") {
      try {
        array.push(command.data.toJSON());
        command_int++;
        console.log(`${notice} ${file} が追加されました。`);
      } catch (error) {
        console.log(
          `\u001b[31m${notice} ${file} はエラーにより追加されませんでした。\nエラー内容\n ${error}\u001b[0m`
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
      console.log(`\u001b[32m===${folder} commands===\u001b[0m`);
      const commandFiles = fs
        .readdirSync(`./commands/${folder}`)
        .filter((file) => file.endsWith(".js"));
      for (const file of commandFiles) {
        const command = require(`../commands/${folder}/${file}`);
        if (command.adminGuildOnly) {
          cmdToArray(adminGuildCommands, command, file, "[Admin]");
          continue;
        }
        if (command.onlyCommand) continue;
        cmdToArray(globalCommands, command, file, "[Global]");
      }
      console.log(`\u001b[32m===${folder} added===\u001b[0m`);
    }

    (async () => {
      try {
        console.log(`${command_int}個のスラッシュコマンドを登録/再登録します…`);

        //Admin
        putToDiscord(adminGuildCommands, testGuild);
        console.log("管理コマンドを正常に登録しました。");

        //Global
        putToDiscord(globalCommands);
        console.log("グローバルコマンドを正常に登録しました。");

        console.log("全てのスラッシュコマンドを正常に登録しました！");
      } catch (error) {
        console.error(error);
      }
    })();
  },
  async handling(client, fs, Collection, config) {
    // コマンドハンドリング
    client.commands = new Collection();
    const commandFolders = fs.readdirSync(`./commands`);
    for (const folder of commandFolders) {
      console.log(`\u001b[32m===${folder} commands===\u001b[0m`);
      const commandFiles = fs
        .readdirSync(`./commands/${folder}`)
        .filter((file) => file.endsWith(".js"));
      for (const file of commandFiles) {
        console.log(`dir:${folder},file:${file}`);
        const command = require(`../commands/${folder}/${file}`);
        try {
          client.commands.set(command.data.name, command);
          console.log(`${command.data.name} がロードされました。`);
        } catch (error) {
          console.log(
            `\u001b[31m${command.data.name} はエラーによりロードされませんでした。\nエラー内容\n ${error}\u001b[0m`
          );
        }
      }
      console.log(`\u001b[32m===${folder} loaded===\u001b[0m`);
    }
  },
};