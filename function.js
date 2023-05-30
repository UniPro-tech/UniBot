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
        return_str = { "year": year, "month": month, "date": date, "hour": hour, "minute": minute, "second": second };
    }
    return return_str;
}
const fs = require('fs');
exports.loging = (data, cmdName) => {
    if (data == "err") {
        throw data;
    }
    // 書き込み
    fs.writeFile(`log/${cmdName}.log`, data, (err) => {
        if (err) throw err;
        console.log('正常に書き込みが完了しました');
    });
}