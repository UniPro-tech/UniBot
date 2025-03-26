/**
 * Converts a timestamp to JST (Japan Standard Time) timestamp.
 * @param {number} timestamp - The timestamp to convert.
 * @returns {Date} - The converted JST timestamp.
 */
const timeToJSTTimestamp = (timestamp: number | string | Date): Date => {
  var dt = new Date(); //Date オブジェクトを作成
  var tz = dt.getTimezoneOffset(); //サーバーで設定されているタイムゾーンの世界標準時からの時差（分）
  tz = (tz + 540) * 60 * 1000; //日本時間との時差(9時間=540分)を計算し、ミリ秒単位に変換

  dt = new Date(timestamp.toString() + tz.toString()); //時差を調整した上でタイムスタンプ値を Date オブジェクトに変換
  return dt;
};

/**
 * Converts a JST (Japan Standard Time) timestamp to a formatted date string or an object with individual date components.
 * @param {number} timestamp - The JST timestamp to convert.
 * @param {boolean} [format=false] - Determines whether to return a formatted date string or an object with individual date components. Default is false.
 * @returns {string|Object} - The formatted date string or an object with individual date components.
 */
export const timeToJST = (
  timestamp: number | string | Date,
  format = false
): Object => {
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
