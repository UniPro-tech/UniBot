import { StreamingDataManager } from "@/lib/dataUtils";
import {
  ChatInputCommandInteraction,
  GuildMember,
  MessageFlags,
  SlashCommandSubcommandBuilder,
  TextChannel,
  VoiceChannel,
} from "discord.js";

export const data = new SlashCommandSubcommandBuilder()
  .setName("enable")
  .setDescription("配信モードの有効化")
  .addChannelOption((option) =>
    option
      .setName("channel")
      .setDescription("配信モードを有効にするチャンネル(デフォルトは現在のチャンネル)")
      .setRequired(false)
  )
  .addStringOption((option) =>
    option.setName("ttsChannel").setDescription("読み上げているチャンネル").setRequired(false)
  );

export const execute = async (interaction: ChatInputCommandInteraction) => {
  if (!interaction.guild) {
    await interaction.reply({
      content: "このコマンドはサーバー内でのみ使用できます。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  let channel = interaction.options.getChannel("channel");
  if (!channel) {
    const member = interaction.member as GuildMember;
    const voiceChannel = member.voice.channel;
    if (!voiceChannel) {
      await interaction.reply({
        content: "チャンネルが指定されておらず、あなたはボイスチャンネルに参加していません。",
        flags: [MessageFlags.Ephemeral],
      });
      return;
    }
    channel = voiceChannel;
  }
  const ttsChannel = interaction.options.getChannel("ttsChannel");
  if (ttsChannel && ttsChannel instanceof TextChannel === false) {
    await interaction.reply({
      content: "読み上げチャンネルにはテキストチャンネルを指定してください。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }
  if (channel instanceof VoiceChannel === false) {
    await interaction.reply({
      content: "配信モードを有効にするチャンネルにはボイスチャンネルを指定してください。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const data = new StreamingDataManager(interaction.guild!.id, channel.id, ttsChannel?.id);

  await data.save();

  await interaction.reply({
    content: `配信モードがチャンネル ${channel} で有効になりました。\n読み上げチャンネル: ${
      ttsChannel || "未指定"
    }`,
    flags: [MessageFlags.Ephemeral],
  });
};
