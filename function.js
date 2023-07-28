const URI_base = process.env.LOGING_URI_BASE;

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
    const month = dt.getMonth() + 1
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

/*LogingAPIにポストして、DBにプッシュする。
*post_data = ポストするためのjson(object)
*api_name = URI
*/
const fs = require('fs');
const https = require('https');
exports.loging = (post_data, api_name) => {
    const URI = `https://${URI_base}/${api_name}`;
    if (post_data == "err") {
        throw post_data;
    }
    // 書き込み
    const options = {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    };
    const request = https.request(URI, options);
    request.write(post_data);
    request.end();
}
exports.readLog = (api_name) => {
    const URI = `https://${URI_base}/${api_name}`;
    https.get(URI, function (res) {
        res.on('data', function (chunk) {
            data.push(chunk);
        }).on('end', function () {

            var events = Buffer.concat(data);
            var r = JSON.parse(events);

            console.log(r);
            return r;

        });
    });
}
