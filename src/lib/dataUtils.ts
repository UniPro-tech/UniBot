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

export default {
  writeConfig,
  readConfig,
};
