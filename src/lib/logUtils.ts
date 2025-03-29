import fs from "fs";
import timeUtils from "@/lib/timeUtils";
import path from "path";
/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
export const write = async (post_data: Object, api_name: string) => {
  const URI = path.resolve(__dirname, `../log/${api_name}`);
  try {
    if (!fs.existsSync(URI)) {
      console.log(
        `[${timeUtils.timeToJST(Date.now(), true)} info]Create directory ${URI}`
      );
      await fs.promises.mkdir(URI, { recursive: true });
    }
    const data = JSON.stringify(post_data);
    fs.writeFile(`${URI}.log`, data, (err: any) => {
      if (err) {
        console.error(
          `\u001b[31m[${timeUtils.timeToJST(
            Date.now(),
            true
          )} error]An Error Occured.\nDatails:\n${err}\u001b[0m`
        );
        throw err;
      } else {
        console.log(
          `[${timeUtils.timeToJST(
            Date.now(),
            true
          )} info]Write data to ${URI}.log`
        );
      }
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
export const read = async (api_name: string) => {
  const URI = path.resolve(__dirname, `../log/${api_name}`);
  if (!fs.existsSync(URI + ".log")) return null;
  try {
    const jsonString = fs.readFileSync(URI + ".log");
    const data = JSON.parse(jsonString.toString());
    return data;
  } catch (error) {
    console.error(
      `\u001b[31m[${timeUtils.timeToJST(
        Date.now(),
        true
      )} error]An Error Occured.\nDatails:\n${error}\u001b[0m`
    );
  }
};

export default {
  write,
  read,
};
