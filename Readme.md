# UntitledBot

DiscordのJavaScriptで書かれた有益なBot

## Getting Started

You can run this Bot on Replit!!

### Replit

First,you've to import to Replit.

And,you have to install [this repogitory](https://github.com/yuito-it/UntitledBot-Loging).

### Config

Write Replit Seacret.
```json
{
  "LOGING_URI_BASE": "Your Loging System's URL",
  "DISCORD_TOKEN": "Discord TOKEN",
}
```
And rename config.exsample.js config.js.
Write config.js.

### GAS
```javascript
var siteURL = "your replit url"
var STAT_OK = "OK";
var STAT_NG = "NG";

function myFunction() {

  try {

    // URLをフェッチ - muteHttpExceptions:trueの場合、HTTPエラーの際に例外をスローしない
    var response = UrlFetchApp.fetch(siteURL, { muteHttpExceptions:true });
    // レスポンスコード
    var code = response.getResponseCode();

    // スクリプトプロパティ
    var scriptProperties = PropertiesService.getScriptProperties();
    // スクリプトプロパティからサイトの死活状態を取得
    var siteStat =  scriptProperties.getProperty('SITE_STAT');

    // レスポンスコード 200をチェックする
    if(code == 200) {
      if(siteStat != STAT_OK )　{
        // サイト状況OKをスクリプトプロパティにセット
        setSiteStat(STAT_OK);
        send("access OK\n" + "code: " + code + "\n");
        console.log('OK code:'+code);
      } 
    } else {
      if(siteStat != STAT_NG )　{
        // サイト状況NGをスクリプトプロパティにセット
        setSiteStat(STAT_NG);
        console.log('NG code:'+code);
        send("access NG \n" + "code: " + code + "\n");
      }
    }

  } catch(err) {
    // catch : DNSエラーなどでURLをfetch出来ないとき

    // サイト状況NGをスクリプトプロパティにセット
    setSiteStat(STAT_NG);
    send("access NG \n" + err  + "\n");
    console.log('err');
  }

}
function send(text) {
  // discord側で作成したボットのウェブフックURL
  const discordWebHookURL = "your log ch webhook url";

  // 投稿するチャット内容と設定
  Logger.log(text);
  const message = {
    "content": text, // チャット本文
    "tts": false  // ロボットによる読み上げ機能を無効化
  }

  const param = {
    "method": "POST",
    "headers": { 'Content-type': "application/json" },
    "payload": JSON.stringify(message)
  }

  UrlFetchApp.fetch(discordWebHookURL, param);
}

// スクリプトプロパティに死活状態をセット
// bStat - OK / NG
function setSiteStat(bStat){

  // スクリプトプロパティをセット
  var scriptProperties = PropertiesService.getScriptProperties();
  scriptProperties.setProperty('SITE_STAT', bStat);

}
```
Set trigger run often 3 min.

## Built With

* [Discord.js](https://discordjs.dev/#/) - The flame work.
* [Deta.space](https://deta.space/) - Loging system.

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

For the versions available, see the [tags on this repository](https://github.com/yuito-it/UntitledBot/tags). 

## Authors

* **Yuito** - [yuito-it](https://github.com/yuito-it)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
