import { PrismaClient } from "@prisma/client";
import { ALStorage, loggingSystem } from "..";

export const prismaClient = new PrismaClient();

/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const writeConfig = async (postData: Object, key: string) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "writeConfig" });
  try {
    await prismaClient.config.upsert({
      where: { key },
      update: { value: JSON.stringify(postData) },
      create: { key, value: JSON.stringify(postData) },
    });
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while writing config"
    );
  }
};

/**
 * Reads a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const readConfig = async (key: string) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "readConfig" });
  try {
    const config = await prismaClient.config.findUnique({
      where: { key },
    });
    return config ? JSON.parse(config.value) : null;
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while reading config"
    );
    return null;
  }
};

export enum SelectedDataType {
  Message = "Message",
  User = "User",
}

export type SelectedData = {
  id?: string;
  user: string;
  type: SelectedDataType;
  data: string;
};

/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const writeSelected = async (data: SelectedData): Promise<void> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "writeSelected" });
  try {
    await prismaClient.selectedData.create({
      data: {
        user: data.user,
        type: data.type,
        data: JSON.stringify(data.data),
      },
    });
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while writing selected data"
    );
  }
};

/**
 * Reads a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const readSelected = async (
  user?: string,
  type?: SelectedDataType,
  data?: string
): Promise<SelectedData | null> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "readSelected" });
  try {
    const selectedData = await prismaClient.selectedData.findFirst({
      where: { user, type, data },
    });
    return (selectedData as SelectedData) || null;
  } catch (error) {
    // TODO: ここでThrow Errorを投げて、呼び出し元でcatchするようにする
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while reading selected data"
    );
    return null;
  }
};

export const writeTtsConnection = async (
  guild: string,
  textChannel: string[],
  voiceChannel: string
): Promise<void> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "writeTtsConnection" });
  try {
    await prismaClient.ttsConnection.upsert({
      where: { guild },
      update: { voiceChannel, textChannel },
      create: { guild, textChannel, voiceChannel },
    });
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while writing TTS connection"
    );
  }
};

export const readTtsConnection = async (
  guild: string,
  textChannel?: string,
  voiceChannel?: string
): Promise<{
  id: string;
  guild: string;
  created_at: Date;
  updated_at: Date;
  textChannel: string[];
  voiceChannel: string;
} | null> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "readTtsConnection" });
  try {
    const connection = await prismaClient.ttsConnection.findFirst({
      where: {
        guild,
        ...(textChannel ? { textChannel: { has: textChannel } } : {}),
        ...(voiceChannel ? { voiceChannel } : {}),
      },
    });
    return connection;
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while reading TTS connection"
    );
    return null;
  }
};

export const writeTtsPreference = async (user: string, key: string, value: object) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "writeTtsPreference" });
  try {
    await prismaClient.ttsPreference.upsert({
      where: { user_key: { user, key } },
      update: { value: JSON.stringify(value) },
      create: { user, key, value: JSON.stringify(value) },
    });
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while writing TTS preference"
    );
  }
};

export const readTtsPreference = async (user: string, key: string): Promise<any | null> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "readTtsPreference" });
  try {
    const preference = await prismaClient.ttsPreference.findFirst({
      where: { user, key },
    });
    return preference ? JSON.parse(preference.value) : null;
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while reading TTS preference"
    );
    return null;
  }
};

export const writeTtsDictionary = async (
  user: string,
  guild: string,
  word: string,
  definition: string,
  caseSensitive: boolean = false
) => {
  try {
    await prismaClient.ttsDictionary.create({
      data: {
        user,
        guild,
        word,
        definition,
        case_sensitive: caseSensitive,
      },
    });
  } catch (error) {
    throw error;
  }
};

export const readTtsDictionary = async (
  user: string,
  guild: string,
  word: string
): Promise<{
  id: string;
  user: string;
  guild: string;
  word: string;
  definition: string;
} | null> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "readTtsDictionary" });
  try {
    const entry = await prismaClient.ttsDictionary.findFirst({
      where: { user, guild, word },
      select: {
        id: true,
        user: true,
        guild: true,
        word: true,
        definition: true,
        case_sensitive: true,
        created_at: true,
        updated_at: true,
      },
      orderBy: [{ created_at: "asc" }, { updated_at: "asc" }],
    });
    return entry;
  } catch (error) {
    logger.error(
      {
        error,
        stack_trace: error instanceof Error ? error.stack : undefined,
      },
      error instanceof Error ? error.message : "An error occurred while reading TTS dictionary"
    );
    return null;
  }
};

export const listTtsDictionary = async (
  guild: string,
  user?: string
): Promise<
  Array<{
    id: string;
    user: string;
    guild: string;
    word: string;
    definition: string;
    caseSensitive: boolean;
    createdAt: Date;
    updatedAt: Date;
  }>
> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "listTtsDictionary" });
  try {
    const entries = await prismaClient.ttsDictionary.findMany({
      where: { user, guild },
      select: {
        id: true,
        user: true,
        guild: true,
        word: true,
        definition: true,
        case_sensitive: true,
        created_at: true,
        updated_at: true,
      },
      orderBy: [{ created_at: "asc" }, { updated_at: "asc" }],
    });
    const formattedEntries = entries.map((entry) => ({
      id: entry.id,
      user: entry.user,
      guild: entry.guild,
      word: entry.word,
      definition: entry.definition,
      caseSensitive: entry.case_sensitive,
      createdAt: entry.created_at,
      updatedAt: entry.updated_at,
    }));
    return formattedEntries;
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while listing TTS dictionary"
    );
    return [];
  }
};

export const removeTtsDictionary = async (
  guild?: string,
  word?: string,
  user?: string,
  id?: string
): Promise<void> => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "removeTtsDictionary" });
  try {
    await prismaClient.ttsDictionary.deleteMany({
      where: { user, guild, word, id },
    });
  } catch (error) {
    logger.error(
      { error, stack_trace: error instanceof Error ? error.stack : undefined },
      error instanceof Error ? error.message : "An error occurred while removing TTS dictionary"
    );
  }
};

type ServerConfigChannelValueType = {
  channel: string;
  value: any;
};

type ServerConfigType<T extends "channel" | "global"> = {
  scope: T;
  data: T extends "channel" ? ServerConfigChannelValueType[] : any;
};
export class ServerDataManager {
  private serverId: string;

  constructor(serverId: string) {
    this.serverId = serverId;
  }

  async getConfig(key: string, channel?: string): Promise<any | null> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "ServerDataManager.getConfig" });
    try {
      const config = await prismaClient.serverConfig.findFirst({
        where: { guild: this.serverId, key },
      });
      if (config) {
        const parsedConfig: ServerConfigType<"channel" | "global"> = JSON.parse(config.value);
        if (parsedConfig.scope === "channel" && channel) {
          const channelData = parsedConfig.data.find(
            (item: ServerConfigChannelValueType) => item.channel === channel
          );
          return channelData ? channelData.value : null;
        } else if (parsedConfig.scope === "global") {
          return parsedConfig.data;
        } else {
          return null;
        }
      }
      return null;
    } catch (error) {
      logger.error(
        { error, stack_trace: error instanceof Error ? error.stack : undefined },
        error instanceof Error ? error.message : "An error occurred while reading server config"
      );
      return null;
    }
  }

  async setConfig(key: string, value: any, channel?: string): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "ServerDataManager.setConfig" });
    let data: ServerConfigType<"channel" | "global">;
    // NOTE: `getConfig` returns the value for a given key/channel (not the full stored object)
    // so we must read the raw database record and parse the stored JSON to get the full shape
    try {
      const existingRecord = await prismaClient.serverConfig.findFirst({
        where: { guild: this.serverId, key },
      });

      if (existingRecord && existingRecord.value) {
        // parse stored structure { scope, data }
        try {
          data = JSON.parse(existingRecord.value) as ServerConfigType<"channel" | "global">;
        } catch (err) {
          // If parsing fails, fallback to creating a new structure based on presence of channel
          logger.warn(
            { err },
            "Failed to parse existing serverConfig.value, overwriting with new structure"
          );
          if (channel) {
            data = { scope: "channel", data: [{ channel, value }] };
          } else {
            data = { scope: "global", data: value };
          }
        }
        // If we have a channel-scoped config, update or append entry
        if (data && data.scope === "channel" && channel) {
          const channelDataIndex = (data.data as ServerConfigChannelValueType[]).findIndex(
            (item: ServerConfigChannelValueType) => item.channel === channel
          );
          if (channelDataIndex !== -1) {
            (data.data as ServerConfigChannelValueType[])[channelDataIndex].value = value;
          } else {
            (data.data as ServerConfigChannelValueType[]).push({ channel, value });
          }
        } else if (data && data.scope === "global" && !channel) {
          // replace global value
          data.data = value;
        } else if (channel) {
          // existing record exists but scope mismatch or not channel-scoped: create channel-scoped structure
          data = { scope: "channel", data: [{ channel, value }] };
        } else {
          data = { scope: "global", data: value };
        }
      } else {
        // no existing record: create new structure
        if (channel) {
          data = { scope: "channel", data: [{ channel, value }] };
        } else {
          data = { scope: "global", data: value };
        }
      }
    } catch (err) {
      // If DB read fails, log and fallback to constructing new structure
      logger.error(
        { err, stack_trace: err instanceof Error ? err.stack : undefined },
        err instanceof Error
          ? err.message
          : "An error occurred while reading existing server config"
      );
      if (channel) {
        data = { scope: "channel", data: [{ channel, value }] };
      } else {
        data = { scope: "global", data: value };
      }
    }
    try {
      await prismaClient.serverConfig.upsert({
        where: { guild_key: { guild: this.serverId, key } },
        update: { value: JSON.stringify(data) },
        create: { guild: this.serverId, key, value: JSON.stringify(data) },
      });
    } catch (error) {
      logger.error(
        { error, stack_trace: error instanceof Error ? error.stack : undefined },
        error instanceof Error ? error.message : "An error occurred while writing server config"
      );
    }
  }

  /**
   * 指定したキーに紐づくサーバー設定を削除または更新します。
   *
   * - channel が指定されている場合:
   *   - データベースから当該キーのレコードを取得し、保存されている JSON を ServerConfigType<"channel" | "global"> としてパースします。
   *   - パース結果の scope が "channel" の場合、ServerConfigChannelValueType の配列から
   *     指定された channel と一致するエントリを除外して値を更新します（該当エントリがなければ変更は行われません）。
   * - channel が指定されていない場合:
   *   - 当該サーバー（this.serverId）かつキーに該当する全レコードを削除します。
   *
   * 失敗時は内部でエラーをキャッチしてロギングを行います（例外は再送出されません）。
   *
   * @param key - 削除または更新対象の設定キー
   * @param channel - オプション。チャンネル単位の設定からそのチャンネルのエントリのみを削除する場合に指定するチャンネル ID
   * @returns Promise<void> - 操作完了を示す Promise。エラーは内部でログ出力されるが、呼び出し元にはスローされません。
   *
   * @example
   * // キーに紐づく全レコードを削除
   * await this.deleteConfig("welcome_message");
   *
   * // チャンネル単位の設定から特定チャンネルのエントリのみを削除
   * await this.deleteConfig("welcome_channels", "123456789012345678");
   */
  async deleteConfig(key: string, channel?: string): Promise<void> {
    const ctx = ALStorage.getStore();
    const logger = loggingSystem.getLogger({ ...ctx, function: "ServerDataManager.deleteConfig" });
    try {
      if (channel) {
        const existingRecord = await prismaClient.serverConfig.findFirst({
          where: { guild: this.serverId, key },
        });
        if (existingRecord) {
          let parsedConfig: ServerConfigType<"channel" | "global">;
          try {
            parsedConfig = JSON.parse(existingRecord.value);
          } catch (parseError) {
            logger.error(
              {
                error: parseError,
                stack_trace: parseError instanceof Error ? parseError.stack : undefined,
              },
              parseError instanceof Error
                ? parseError.message
                : "Failed to parse server config JSON in deleteConfig"
            );
            return;
          }
          if (parsedConfig.scope === "channel") {
            parsedConfig.data = parsedConfig.data.filter(
              (item: ServerConfigChannelValueType) => item.channel !== channel
            );
            if (parsedConfig.data.length === 0) {
              await prismaClient.serverConfig.deleteMany({
                where: { guild: this.serverId, key },
              });
            } else {
              await prismaClient.serverConfig.update({
                where: { guild_key: { guild: this.serverId, key } },
                data: { value: JSON.stringify(parsedConfig) },
              });
            }
          }
        }
      } else {
        await prismaClient.serverConfig.deleteMany({
          where: { guild: this.serverId, key },
        });
      }
    } catch (error) {
      logger.error(
        { error, stack_trace: error instanceof Error ? error.stack : undefined },
        error instanceof Error ? error.message : "An error occurred while deleting server config"
      );
    }
  }
}

export default {
  writeConfig,
  readConfig,
  writeSelected,
  readSelected,
  writeTtsConnection,
  readTtsConnection,
  writeTtsPreference,
  readTtsPreference,
  writeTtsDictionary,
  readTtsDictionary,
};
