const { SlashCommandBuilder } = require("discord.js");
const Discord = require("discord.js");

module.exports = {
    data: new SlashCommandBuilder()
        .setName("about")
        .setDescription(
            "アイコンのURLを取得します。"
        ),
    async execute(i, client) {
        let size;
        client.shard.fetchClientValues('guilds.cache.size')
            .then(result => {
                size = result;
            });
        const embed = new Discord.EmbedBuilder()
            .setColor(0x0099FF)
            .setTitle(`About ${client.conf.productname}`)
            //.setURL('https://discord.js.org/')
            .setAuthor(client.conf.author)
            .setDescription(client.conf.description)
            .setThumbnail(client.user.displayAvatarURL({ dynamic: true }))
            .addFields(
                { name: 'Version', value: client.conf.version },
                { name: 'Author', value: client.conf.author.name },
                { name: 'Guild Count', value: `${size}` }
            )
            .setTimestamp();
        i.reply({ embeds: [embed] });
    }
}