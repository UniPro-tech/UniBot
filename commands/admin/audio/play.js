const { SlashCommandSubcommandBuilder, EmbedBuilder } = require("discord.js");
const { useMainPlayer, ExtractorExecutionContext } = require('discord-player');

module.exports = {
    data: new SlashCommandSubcommandBuilder()
        .setName("play")
        .setDescription("指定されたオーディオを再生します。")
        .addStringOption(option => option
            .setName("query")
            .setDescription("再生するオーディオを指定。")),
    async execute(i, client) {
        await i.deferReply();
        const player = useMainPlayer();
        const channel = i.member.voice.channel;
        if (!channel) {
            const embed = new EmbedBuilder()
                .setAuthor({ name: i.user.username, iconURL: i.user.displayAvatarURL({ dynamic: true }) })
                .setColor(client.conf.color.e)
                .setTitle("エラー")
                .setDescription(`VC未参加です。
            VCに参加してください。`)
                .setFooter(client.conf.footer)
                .setTimestamp();
            return i.followUp({ embeds: [embed] }); // make sure we have a voice channel
        }
        const query = i.options.getString('query', true); // we need input/query to play

        try {
            const extcont = new ExtractorExecutionContext(player);
            player.extractors = extcont;
            if (~query.indexOf("spotify")) {
                await player.extractors.loadDefault((ext) => 'SpotifyExtractor');
            }
            else if (~query.indexOf("youtube") || ~query.indexOf("youtu.be")) {
                await player.extractors.loadDefault((ext) => 'YouTubeExtractor');
            }
            else if (~query.indexOf("music.apple")) {
                await player.extractors.loadDefault((ext) => "AppleMusicExtractor");
            }
            else if (~query.indexOf("soundcloud")) {
                await player.extractors.loadDefault((ext) => "SoundCloudExtractor");
            } else {
                await player.extractors.loadDefault((ext) => "YoutubeExtractor");
            }
            const { track } = await player.play(channel, query, {
                /*nodeOptions: {
                    // nodeOptions are the options for guild node (aka your queue in simple word)
                    metadata: i // we can access this metadata object using queue.metadata later on
                }*/
            });
            const embed = new EmbedBuilder()
                .setTitle(`**${track.title}** を再生`)
                .setURL(track.url)
                .setAuthor({ name: i.user.username, iconURL: i.user.displayAvatarURL({ dynamic: true }) })
                .setColor(client.conf.color.s)
                .setDescription(track.description)
                .setThumbnail(track.thumbnail)
                .setFooter(client.conf.footer)
                .setTimestamp();
            return i.followUp({ embeds: [embed] });
        } catch (e) {
            // let's return error if something failed
            console.log(e);
            return i.followUp(`エラー: ${e}`);
        }
    }
}