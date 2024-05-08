const { URI_base } = require("./config");
const fs = require("fs");

//タイムスタンプをJSTタイムスタンプに変換
function timeToJSTTimestamp(timestamp) {
  var dt = new Date(); //Date オブジェクトを作成
  var tz = dt.getTimezoneOffset(); //サーバーで設定されているタイムゾーンの世界標準時からの時差（分）
  tz = (tz + 540) * 60 * 1000; //日本時間との時差(9時間=540分)を計算し、ミリ秒単位に変換

  dt = new Date(timestamp + tz); //時差を調整した上でタイムスタンプ値を Date オブジェクトに変換
  return dt;
}
exports.timeToJSTTimestamp = timeToJSTTimestamp;

//JSTタイムスタンプから日付
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

exports.loging = async (post_data, api_name) => {
  const URI = `./log/${api_name}/`;
  try {
    if (!fs.existsSync(URI)) {
      fs.promises.mkdir(URI, { recursive: true });
    }
    const data = JSON.stringify(post_data);
    await fs.writeFile(`${URI}.log`, data, (err) => {
      // 書き出しに失敗した場合
      if (err) {
        console.log("エラーが発生しました。" + err);
        throw err;
      }
      // 書き出しに成功した場合
      else {
        console.log("ファイルが正常に書き出しされました");
      }
    });
  } catch (e) {
    console.log(e);
  }
};
