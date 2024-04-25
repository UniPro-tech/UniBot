const { SlashCommandBuilder, Activity } = require("discord.js");
const Discord = require("discord.js");
module.exports = {
    guildOnly: false, // サーバー専用コマンドかどうか
    adminGuildOnly: true,
    data: new SlashCommandBuilder() // スラッシュコマンド登録のため
        .setName("maintenance")
        .setDescription("メンテモード")
        .addStringOption(option =>
            option.setName('statusdiscription')
                .setDescription('プレイ中に表示するやつ(Stream、WatchingはURLでもok)'))
        .addStringOption(option =>
            option.setName('status')
                .setDescription('すてーたす')
                .setChoices(
                    { name: 'オンライン', value: 'online' },
                    { name: '取り込み中', value: 'dnd' },
                    { name: 'スリープ', value: 'idle' },
                    //{ name: 'スマホでオンライン', value: 'Discord Android' },
                    { name: 'オンライン隠し', value: 'invisible' }
                ))
        .addStringOption(option =>
            option.setName('activitytype')
                .setDescription('あくてぃびてぃ')
                .setChoices(
                    { name: '視聴中', value: 'WATCHING' },
                    { name: 'プレイ中', value: 'PLAYING' },
                    { name: '参戦中', value: 'COMPETING' },
                    { name: '再生中(聞く)', value: 'LISTENING' },
                    { name: '配信中', value: 'STREAMING' },
                    { name: 'カスタム', value: 'CUSTOM' })),

    async execute(i, client, command) {
        try {
            const onoff = i.options.getString('enable');
            if (onoff == 'on') {
                const status = i.options.getString('status');
                const statusDescription = i.options.getString('statusdiscription');
                const activityType = i.options.getString('activitytype');
                let activityOpt;
                switch (activityType) {
                    case "WATCHING":
                        activityOpt.type = Discord.ActivityType.Watching;
                        break;

                    case "COMPETING":
                        activityOpt.type = Discord.ActivityType.Competing;
                        break

                    case "LISTENING":
                        activityOpt.type = Discord.ActivityType.Listening;
                        break;

                    case "STREAMING":
                        activityOpt.type = Discord.ActivityType.Streaming;
                        break;

                    case "CUSTOM":
                        activityOpt.type = Discord.ActivityType.Custom;

                    default:
                        activityOpt.type = Discord.ActivityType.Playing;
                        break;
                }
                /*
                if (status == 'Discord Android') {
                    client.ws = { properties: { "$os": "Untitled OS", "$browser": "Discord Android", "$device": "Replit Container" } };
                    client.user.setStatus('online');
                } else {*/
                //    client.ws = () => { return { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } }; }
                //}

                client.user.setActivity(statusDescription, activityOpt);
                client.user.setStatus(status);

                const embed = new Discord.EmbedBuilder()
                    .setTitle("ok")
                    .setColor(client.conf.color.s)
                    .setTimestamp();

                i.reply({ embeds: [embed] });
                client.func.loging({ onoff: "on", status: status, playing: statusDescription, type: activityType }, "v1/conf/status");
                return `{ "onoff":"on","status": "${status}", "statusDesc": "${statusDescription}", "type": "${activityType}" }`;
            } else {
                client.shard.fetchClientValues('guilds.cache.size')
                    .then(result => {
                        client.user.setActivity(`Servers: ${result}`);
                    });
                //client.ws = { properties: { "$os": "Untitled OS", "$browser": "Untitled Browser", "$device": "Replit Container" } };
                client.user.setStatus("online");
                const embed = new Discord.EmbedBuilder()
                    .setTitle("ok")
                    .setColor(client.conf.color.s)
                    .setTimestamp();

                i.reply({ embeds: [embed] });
                client.func.loging({ onoff: "off" }, "v1/conf/status");
                return `{ "onoff":"off"}`;
            }
        } catch (e) {
            throw e;
        }

    },
}