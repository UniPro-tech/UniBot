import fs from "fs";
import timeUtils from "@/lib/timeUtils";
import path from "path";

import { PrismaClient } from "@prisma/client";

export const prismaClient = new PrismaClient();

/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const writeConfig = async (post_data: Object, key: string) => {
  try {
    await prismaClient.config.upsert({
      where: { key },
      update: { value: JSON.stringify(post_data) },
      create: { key, value: JSON.stringify(post_data) },
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

export default {
  writeConfig,
  readConfig,
};
