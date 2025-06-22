import { Client, SlashCommandBuilder } from "discord.js";
import { RESTPostAPIChatInputApplicationCommandsJSONBody, Routes } from "discord-api-types/v10";
import { REST } from "@discordjs/rest";
import fs from "fs";
import path from "path";
import { ChatInputCommand } from "@/executors/types/ChatInputCommand";
import { ContextMenuCommand } from "@/executors/types/ContextMenuCommand";

/**
 * 差分比較用にフィールドを絞って整形
 */
function cleanCommandData(cmd: RESTPostAPIChatInputApplicationCommandsJSONBody) {
  const { name, description, options, type, default_member_permissions, dm_permission } = cmd;
  return {
    name,
    description,
    options: options ?? [],
    type: type ?? 1,
    default_member_permissions: default_member_permissions ?? null,
    dm_permission: dm_permission ?? true,
  };
}

function areCommandsEqual(
  a: RESTPostAPIChatInputApplicationCommandsJSONBody,
  b: RESTPostAPIChatInputApplicationCommandsJSONBody
): boolean {
  return JSON.stringify(cleanCommandData(a)) === JSON.stringify(cleanCommandData(b));
}

/**
 * Discord に登録（差分がある場合のみ更新）
 */
async function putToDiscordWithDiffCheck(
  client: Client,
  rest: REST,
  array: RESTPostAPIChatInputApplicationCommandsJSONBody[],
  guild?: string
) {
  const route = guild
    ? Routes.applicationGuildCommands(client.application?.id as string, guild)
    : Routes.applicationCommands(client.application?.id as string);

  const existing = (await rest.get(route)) as RESTPostAPIChatInputApplicationCommandsJSONBody[];

  // 名前をキーにしたMapに変換
  const existingMap = new Map<string, RESTPostAPIChatInputApplicationCommandsJSONBody>();
  existing.forEach((cmd) => existingMap.set(cmd.name, cmd));

  // 長さが違うだけで更新対象にする
  if (array.length !== existing.length) {
    console.info(`[Update] Command count changed: ${existing.length} -> ${array.length}`);
  }

  const shouldUpdate =
    array.length !== existing.length ||
    array.some((cmd) => {
      const existingCmd = existingMap.get(cmd.name);
      if (!existingCmd) {
        // 新規コマンドあり
        return true;
      }
      return !areCommandsEqual(cmd, existingCmd);
    });

  if (!shouldUpdate) {
    console.info(
      `[Skip] No changes detected. Skipping registration for ${guild ?? "global"} commands.`
    );
    return;
  }

  await rest.put(route, { body: array });
  console.info(`[Update] Commands updated for ${guild ?? "global"} scope.`);
}

/**
 * コマンド登録関数（チャット入力＋右クリック含む）
 */
export const registerAllCommands = async (client: Client) => {
  console.info(`\u001b[32m===Pushing All ApplicationCommand Data===\u001b[0m`);
  const config = client.config;
  const token = config.token;
  const testGuild = config.dev.testGuild;
  const rest = new REST({ version: "10" }).setToken(token);

  const globalCommands: RESTPostAPIChatInputApplicationCommandsJSONBody[] = [];
  const adminGuildCommands: RESTPostAPIChatInputApplicationCommandsJSONBody[] = [];

  let commandCount = 0;

  const pushCommand = (
    arr: RESTPostAPIChatInputApplicationCommandsJSONBody[],
    command: ChatInputCommand | ContextMenuCommand,
    file: string,
    typeLabel: string
  ) => {
    try {
      const data = (command.data as SlashCommandBuilder).toJSON();
      arr.push(data);
      commandCount++;
      console.info(`[${typeLabel}] ${file} has been added.`);
    } catch (err) {
      console.error(`[${typeLabel}] Error in ${file}:\n`, err);
    }
  };

  // --- Slash Commands 読み込み
  const slashCommandFolders = fs.readdirSync(
    path.resolve(__dirname, "../executors/chatInputCommands")
  );
  for (const folder of slashCommandFolders) {
    const commandFiles = fs
      .readdirSync(path.resolve(__dirname, `../executors/chatInputCommands/${folder}`))
      .filter((f) => f.endsWith(".js") || (f.endsWith(".ts") && !f.endsWith(".d.ts")));

    for (const file of commandFiles) {
      const command = require(path.resolve(
        __dirname,
        `../executors/chatInputCommands/${folder}/${file}`
      ));
      if (command.adminGuildOnly) {
        pushCommand(adminGuildCommands, command, file, "Admin Slash");
      } else {
        pushCommand(globalCommands, command, file, "Global Slash");
      }
    }
  }

  // --- Context Menu Commands 読み込み
  const contextMenuFiles = fs
    .readdirSync(path.resolve(__dirname, "../executors/messageContextMenuCommands"))
    .filter((f) => f.endsWith(".js") || (f.endsWith(".ts") && !f.endsWith(".d.ts")));

  for (const file of contextMenuFiles) {
    const command = await import(
      path.resolve(__dirname, `../executors/messageContextMenuCommands/${file}`)
    );
    if (command.adminGuildOnly) {
      pushCommand(adminGuildCommands, command, file, "Admin Context");
    } else {
      pushCommand(globalCommands, command, file, "Global Context");
    }
  }

  // --- 実行
  try {
    console.info(`[Init] Registering ${commandCount} total commands...`);
    await putToDiscordWithDiffCheck(client, rest, adminGuildCommands, testGuild);
    await putToDiscordWithDiffCheck(client, rest, globalCommands);
    console.info(`[Done] All commands registered successfully.`);
  } catch (err) {
    console.error("[Register Error]", err);
  }
};
