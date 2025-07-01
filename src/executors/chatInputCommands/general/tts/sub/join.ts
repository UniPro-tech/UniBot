import {
  CommandInteraction,
  EmbedBuilder,
  GuildMember,
  SlashCommandSubcommandBuilder,
  TextChannel,
} from "discord.js";
import {
  createAudioPlayer,
  createAudioResource,
  getVoiceConnection,
  joinVoiceChannel,
} from "@discordjs/voice";
import { writeTtsConnection } from "@/lib/dataUtils";
import { Readable } from "stream";
import { RPC, Query, Generate } from "voicevox.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("join")
  .setDescription("ボイスチャンネルに参加します。");

const createErrorEmbed = (title: string, description: string, color: number) =>
  new EmbedBuilder().setTitle(title).setDescription(description).setColor(color).setTimestamp();

const createSuccessEmbed = (voiceChannelId: string, textChannelId: string, color: number) =>
  new EmbedBuilder()
    .setTitle("TTSボイスチャンネル接続")
    .setDescription("ボイスチャンネルに接続しました。")
    .addFields([
      { name: "ボイスチャンネル名", value: `<#${voiceChannelId}>` },
      { name: "テキストチャンネル名", value: `<#${textChannelId}>` },
    ])
    .setColor(color);

const connectVoiceVox = async () => {
  if (!RPC.rpc) {
    const headers = {
      Authorization: `ApiKey ${process.env.VOICEVOX_API_KEY}`,
    };
    await RPC.connect(process.env.VOICEVOX_API_URL as string, headers);
  }
};

export const execute = async (interaction: CommandInteraction) => {
  await interaction.deferReply();

  const member = interaction.member as GuildMember;
  const voiceChannel = member.voice.channel;
  const config = interaction.client.config;

  if (!voiceChannel) {
    const embed = createErrorEmbed(
      "Error - TTSボイスチャンネル接続失敗",
      "ボイスチャンネルに参加していません。",
      config.color.error
    );
    await interaction.editReply({ embeds: [embed] });
    return;
  }

  if (getVoiceConnection(interaction.guild!.id)) {
    const embed = createErrorEmbed(
      "Error - 既に接続中",
      "既にボイスチャンネルに接続しています。",
      config.color.error
    );
    await interaction.editReply({ embeds: [embed] });
    return;
  }

  const connection = joinVoiceChannel({
    channelId: voiceChannel.id,
    guildId: voiceChannel.guild.id,
    adapterCreator: voiceChannel.guild.voiceAdapterCreator,
  });

  connection.once("ready", async () => {
    const textChannelId = (interaction.channel! as TextChannel).id;
    const embed = createSuccessEmbed(voiceChannel.id, textChannelId, config.color.success);
    await interaction.editReply({ embeds: [embed] });

    const player = createAudioPlayer();
    connection.subscribe(player);

    const text = `${voiceChannel.name}に接続しました。`;
    await connectVoiceVox();
    const query = await Query.getTalkQuery(text, 0);
    const audio = await Generate.generate(0, query);
    const audioStream = Readable.from(audio);
    const resource = createAudioResource(audioStream);
    player.play(resource);

    writeTtsConnection(voiceChannel.guild.id, [interaction.channel?.id as string], voiceChannel.id);
  });
};

export default {
  data,
  execute,
};
