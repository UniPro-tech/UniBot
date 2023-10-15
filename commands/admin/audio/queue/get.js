const { SlashCommandSubcommandBuilder, EmbedBuilder, Client } = require("discord.js");
const { useQueue, Track, GuildQueue } = require("discord-player");

module.exports = {
    data: new SlashCommandSubcommandBuilder()
        .setName("get")
        .setDescription("キューの中身を表示")
        .addNumberOption((option) => option
            .setName("page")
            .setDescription("ページ番号")
        ),
    async execute(i, client) {
        await i.deferReply();

        const queue = useQueue(i.guildId ?? "");

        if (queue) {
            const page = i.options.getNumber("page") ?? 1;
            const tracks = queue.tracks.toArray();
            tracks.unshift(queue.history.currentTrack);
            console.log(tracks.length)

            const limit_fields = 25;
            const pageMax = Math.ceil(tracks.length / limit_fields);

            const embed = new EmbedBuilder()
                .setAuthor({ name: i.user.username, iconURL: i.user.displayAvatarURL({ dynamic: true }) })
                .setColor(client.conf.color.s)
                .setFooter(client.conf.footer)
                .setTitle("現在のキュー")

            const selectField = 25 * (page - 1);
            console.log(selectField);
            console.log(selectField + 25);

            for (let index = selectField; index < selectField + 25 || index + 1 < tracks.length; index++) {
                const num = index - selectField;
                let name;
                if (index==0)name = "Now playing";else name = num.toString;
                const t = tracks[index];
                if (t && t.title) {
                    console.log(`${index}[${t.title} • ${t.author}](${t.url}) (${t.duration})`);
                    await embed.addFields({ name: name, value: `${index}[${t.title} • ${t.author}](${t.url}) (${t.duration})` });
                } else {
                    console.log(`Track at index ${index} is missing 'title' property.`);
                }
            }
            await i.followUp({ embeds: [embed] });
        } else {
            const embed = new EmbedBuilder()
                .setColor(client.conf.color.e)
                .setTitle("キューがありません。")
                .setDescription(`原因はおそらく下のどれかです。
            - 何も再生していない。
            - あなた自身がVCに参加していない。
            これで解決しない場合は、サポートへ連絡してください。`)
                .setTimestamp()
                .setFooter(client.conf.footer);
            return i.followUp({ embeds: [embed] });
        }
    }
}