# UniBot

UniPro内を統治したりしているDiscord Botです。

## ローカルで動かす

### 前提条件

- Bun >=1.2

### 設定

_envを.envとしてコピーし、各種値を変更してください。

### 実行

#### 依存関係のインストール

npm packageのインストールを行います。

```shell
bun install
```

#### 実行

Bunのランタイムで実行します。

```shell
bun src/index.ts
```

## Dockerで動かす

[docker-composeファイル](docker-compose.prod.yaml)を実行します。

## 開発する

ローカルで開発する場合は、以上の手順に従ってください。

コンテナ内で開発したい場合は、[開発用のdocker-composeファイル](docker-compose.dev.yaml)を用いてください。

## Built With

- [Discord.js](https://discordjs.dev/#/) - The flame work.

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

For the versions available, see the [tags on this repository](https://github.com/yuito-it/UntitledBot/tags).

## Authors

- @yuito-it

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

- Hat tip to anyone whose code was used
- Inspiration
- etc
