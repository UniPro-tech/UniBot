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
  .setName("disable")
  .setDescription("配信モードの無効化")
  .addChannelOption((option) =>
    option
      .setName("channel")
      .setDescription("配信モードを無効にするチャンネル(デフォルトは現在のチャンネル)")
      .setRequired(false)
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
  if (channel instanceof VoiceChannel === false) {
    await interaction.reply({
      content: "配信モードを無効にするチャンネルにはボイスチャンネルを指定してください。",
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const data = new StreamingDataManager(interaction.guild!.id, channel.id, undefined);

  await data.delete();

  await interaction.reply({
    content: `配信モードがチャンネル ${channel} で無効になりました。`,
    flags: [MessageFlags.Ephemeral],
  });
};
