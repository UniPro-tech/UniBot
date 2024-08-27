const fs = require("fs");
/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
exports.loging = async (post_data, api_name) => {
  const URI = `${process.cwd()}/log/${api_name}/`;
  try {
    if (!fs.existsSync(URI)) {
      fs.promises.mkdir(URI, { recursive: true });
    }
    const data = JSON.stringify(post_data);
    await fs.writeFile(`${URI}.log`, data, (err) => {
      if (err) {
        console.log("エラーが発生しました。" + err);
        throw err;
      } else {
        console.log("ファイルが正常に書き出しされました");
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
    console.error("エラー:", error.message);
  }
};
