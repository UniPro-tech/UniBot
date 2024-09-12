const fs = require("fs");
const timeUtils = require("./timeUtils.js");
/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
exports.loging = async (post_data, api_name) => {
  const URI = `${process.cwd()}/log/${api_name}`;
  try {
    if (!fs.existsSync(URI)) {
      fs.promises.mkdir(URI, { recursive: true });
    }
    const data = JSON.stringify(post_data);
    fs.writeFile(`${URI}.log`, data, (err) => {
      if (err) {
        console.error(
          `\u001b[31m[${timeUtils.timeToJST(
            Date.now(),
            true
          )} error]An Error Occured.\nDatails:\n${err}\u001b[0m`
        );
        throw err;
      } else {
        console.log(`[${client.func.timeUtils.timeToJST(Date.now(), true)} info]Write data to ${URI}.log`);
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
exports.readLog = async (api_name) => {
  const URI = `${process.cwd()}/log/${api_name}/`;
  try {
    const jsonString = fs.readFileSync(URI + ".log");
    const data = JSON.parse(jsonString);
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
