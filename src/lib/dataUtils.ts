import timeUtils from "@/lib/timeUtils";

import { PrismaClient, selectedData } from "@prisma/client";

export const prismaClient = new PrismaClient();

/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const writeConfig = async (postData: Object, key: string) => {
  try {
    await prismaClient.config.upsert({
      where: { key },
      update: { value: JSON.stringify(postData) },
      create: { key, value: JSON.stringify(postData) },
    });
  } catch (e) {
    console.log(e);
  }
};

/**
 * Reads a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const readConfig = async (key: string) => {
  try {
    const config = await prismaClient.config.findUnique({
      where: { key },
    });
    return config ? JSON.parse(config.value) : null;
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured.\nDatails:\n${error}\u001b[0m`
    );
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
  try {
    await prismaClient.selectedData.create({
      data: {
        user: data.user,
        type: data.type,
        data: JSON.stringify(data.data),
      },
    });
  } catch (e) {
    console.error(e);
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
  try {
    const selectedData = await prismaClient.selectedData.findFirst({
      where: { user, type, data },
    });
    return (selectedData as SelectedData) || null;
  } catch (error) {
    // TODO: ここでThrow Errorを投げて、呼び出し元でcatchするようにする
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured.\nDatails:\n${error}\u001b[0m`
    );
    return null;
  }
};

export const writeTtsConnection = async (
  guild: string,
  textChannel: string[],
  voiceChannel: string
): Promise<void> => {
  try {
    await prismaClient.ttsConnection.upsert({
      where: { guild },
      update: { voiceChannel, textChannel },
      create: { guild, textChannel, voiceChannel },
    });
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
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
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occured.\nDatails:\n${error}\u001b[0m`
    );
    return null;
  }
};

export const writeTtsPreference = async (user: string, key: string, value: object) => {
  try {
    await prismaClient.ttsPreference.upsert({
      where: { user_key: { user, key } },
      update: { value: JSON.stringify(value) },
      create: { user, key, value: JSON.stringify(value) },
    });
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An error occurred.\nDetails:\n${error}\u001b[0m`
    );
  }
};

export const readTtsPreference = async (user: string, key: string): Promise<any | null> => {
  try {
    const preference = await prismaClient.ttsPreference.findFirst({
      where: { user, key },
    });
    return preference ? JSON.parse(preference.value) : null;
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
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
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
    );
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
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
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
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
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
  try {
    await prismaClient.ttsDictionary.deleteMany({
      where: { user, guild, word, id },
    });
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJSTstamp(
        Date.now(),
        true
      )} error]An Error Occurred.\nDetails:\n${error}\u001b[0m`
    );
  }
};

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
