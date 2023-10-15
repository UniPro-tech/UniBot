const { SlashCommandSubcommandBuilder } = require("discord.js");
const { useQueue } = require("discord-player");

module.exports = {
    data: new SlashCommandSubcommandBuilder()
        .setName("skip")
        .setDescription("曲をスキップ"),
    async execute(i, client) {
        await i.deferReply();
        const { member } = i;

        try {
            const queue = useQueue(i.guild.id ?? "");
            // #requireVoiceSession doesn't check current track,
            // only session/player state
            const currentTrack = queue?.currentTrack;
            if (!currentTrack) {
                i.followUp({ content: ` ${member}, 何も再生されていません...` });
                return;
            }
            await queue.node.skip();
            await i.followUp(
                ` ${member}により、 **\`${currentTrack}\`** がスキップされました。`);
        }
        catch (e) {
            i.followUp(` ${member}, エラー: \n\n${e.message}`);
        }
    }
}