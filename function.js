const fs = require("fs");

/**
 * Converts a timestamp to JST (Japan Standard Time) timestamp.
 * @param {number} timestamp - The timestamp to convert.
 * @returns {Date} - The converted JST timestamp.
 */
function timeToJSTTimestamp(timestamp) {
  var dt = new Date(); //Date オブジェクトを作成
  var tz = dt.getTimezoneOffset(); //サーバーで設定されているタイムゾーンの世界標準時からの時差（分）
  tz = (tz + 540) * 60 * 1000; //日本時間との時差(9時間=540分)を計算し、ミリ秒単位に変換

  dt = new Date(timestamp + tz); //時差を調整した上でタイムスタンプ値を Date オブジェクトに変換
  return dt;
}
exports.timeToJSTTimestamp = timeToJSTTimestamp;

/**
 * Converts a JST (Japan Standard Time) timestamp to a formatted date string or an object with individual date components.
 * @param {number} timestamp - The JST timestamp to convert.
 * @param {boolean} [format=false] - Determines whether to return a formatted date string or an object with individual date components. Default is false.
 * @returns {string|Object} - The formatted date string or an object with individual date components.
 */
exports.timeToJST = function (timestamp, format = false) {
  const dt = timeToJSTTimestamp(timestamp);
  const year = dt.getFullYear();
  const month = dt.getMonth() + 1;
  const date = dt.getDate();
  const hour = dt.getHours();
  const minute = dt.getMinutes();
  const second = dt.getSeconds();

  let return_str;
  if (format == true) {
    return_str = `${year}/${month}/${date} ${hour}:${minute}:${second}`;
  } else {
    return_str = {
      year: year,
      month: month,
      date: date,
      hour: hour,
      minute: minute,
      second: second,
    };
  }
  return return_str;
};

/**
 * Reads a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
exports.readLog = async (api_name) => {
  const URI = `./log/${api_name}/`;
  try {
    const jsonString = fs.readFileSync(URI + ".log");
    const data = JSON.parse(jsonString);
    return data;
  } catch (error) {
    console.error("エラー:", error.message);
  }
};

/**
 * Writs a log file for a specific directory.
 * @param {string} api_name - The name of the path.
 * @returns {Promise<Object>} - The parsed log data.
 */
exports.loging = async (post_data, api_name) => {
  const URI = `./log/${api_name}/`;
  try {
    if (!fs.existsSync(URI)) {
      fs.promises.mkdir(URI, { recursive: true });
    }
    const data = JSON.stringify(post_data);
    await fs.writeFile(`${URI}.log`, data, (err) => {
      if (err) {
        console.log("エラーが発生しました。" + err);
        throw err;
      }
      else {
        console.log("ファイルが正常に書き出しされました");
      }
    });
  } catch (e) {
    console.log(e);
  }
};
