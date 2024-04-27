const { URI_base } = require("./config");

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

/*LogingAPIにポストして、DBにプッシュする。
 *post_data = ポストするためのjson(object)
 *api_name = URI
 */
const axios = require("axios");
const https = require("https");
exports.readLog = async (api_name) => {
  const URI = `${URI_base}/${api_name}`;
  try {
    const response = await axios.get(URI);
    const ret = JSON.stringify(response.data);
    console.log("レスポンス:", response.data);
    return JSON.parse(ret);
  } catch (error) {
    console.error("エラー:", error.message);
  }
};
const http = require("http");

exports.loging = async (post_data, api_name) => {
  const url = `${URI_base}/${api_name}`;
  console.log("URI:", url);
  try {
    const response = await axios.post(url, post_data);

    console.log("Response:", response.data);
  } catch (error) {
    console.error("Error:", error.message);
  }
};
