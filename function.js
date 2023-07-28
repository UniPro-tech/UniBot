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
const http = require('http')

if (require.main === module) {
    main()
}

exports.loging = async(post_data, api_name) => {
    try {
        const res = await new Promise((resolve, reject) => {
            try {
                const port = process.env.PORT || '8080'
                const url = `http://${URI_base}/${api_name}`
                const content = post_data
                const req = http.request(url, { // <1>
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Content-Length': '' + content.length,
                        'X-Header': 'X-Header',
                    },
                })

                req.on('response', resolve) // <2>
                req.on('error', reject) // <3>
                req.write(content) // <4>
                req.end() // <5>
            } catch (err) {
                reject(err)
            }
        })

        const chunks = await new Promise((resolve, reject) => {
            try {
                const chunks = []

                res.on('data', (chunk) => chunks.push(chunk)) // <6>
                res.on('end', () => resolve(chunks)) // <7>
                res.on('error', reject) // <8>
            } catch (err) {
                reject(err)
            }
        })

        const buffer = Buffer.concat(chunks) // <9>
        const body = JSON.parse(buffer)

        console.info(JSON.stringify(body, null, 2))
        return body;
    } catch (err) {
        console.error(err)
    }
}