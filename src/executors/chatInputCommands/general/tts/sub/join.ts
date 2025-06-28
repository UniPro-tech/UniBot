import {
  CommandInteraction,
  EmbedBuilder,
  GuildMember,
  SlashCommandSubcommandBuilder,
  TextChannel,
} from "discord.js";
import { createAudioPlayer, createAudioResource, joinVoiceChannel } from "@discordjs/voice";
import { writeTtsConnection } from "@/lib/dataUtils";
import { Readable } from "stream";
import { RPC, Query, Generate } from "voicevox.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("join")
  .setDescription("ボイスチャンネルに参加します。");
export const execute = async (interaction: CommandInteraction) => {
  await interaction.deferReply();
  const voiceChannel = (interaction.member as GuildMember).voice.channel;
  if (!voiceChannel) {
    await interaction.followUp("ボイスチャンネルに参加していません。");
    return;
  }
  const connection = await joinVoiceChannel({
    channelId: voiceChannel.id,
    guildId: voiceChannel.guild.id,
    adapterCreator: voiceChannel.guild.voiceAdapterCreator,
  });
  connection.once("ready", async () => {
    console.info("Connected to voice channel");
    const embed = new EmbedBuilder()
      .setTitle("TTSボイスチャンネル接続")
      .setDescription(`ボイスチャンネルに接続しました。`)
      .addFields([
        {
          name: "ボイスチャンネル名",
          value: `<#${voiceChannel.id}>`,
        },
        {
          name: "テキストチャンネル名",
          value: `<#${(interaction.channel! as TextChannel).id}>`,
        },
      ])
      .setColor(interaction.client.config.color.success);
    await interaction.editReply({ embeds: [embed] });
    writeTtsConnection(voiceChannel.guild.id, [interaction.channel?.id as string], voiceChannel.id);
    const player = createAudioPlayer();
    connection.subscribe(player);
    const text = `${voiceChannel.name}に接続しました。`;
    if (!RPC.rpc) {
      const headers = {
        Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
      };
      await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
    }
    const query = await Query.getTalkQuery(text, 0);
    const audio = await Generate.generate(0, query);
    const audioStream = Readable.from(audio);
    const resource = createAudioResource(audioStream);
    player.play(resource);
  });
};

export default {
  data,
  execute,
};
