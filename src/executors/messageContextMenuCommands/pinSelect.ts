import {
  ContextMenuCommandBuilder,
  ApplicationCommandType,
  MessageContextMenuCommandInteraction,
  MessageFlags,
  EmbedBuilder,
  ChannelType,
  PartialGroupDMChannel,
} from "discord.js";
import { ServerDataManager } from "@/lib/dataUtils.js";
import config from "@/config.js";

export const name = "Pinするメッセージを選択";

export const data = new ContextMenuCommandBuilder()
  .setName(name)
  .setType(ApplicationCommandType.Message);

export const execute = async (interaction: MessageContextMenuCommandInteraction) => {
  const isAdmin =
    interaction.memberPermissions?.has("ManageMessages") ||
    interaction.channel?.type === ChannelType.DM;
  const targetMsg = interaction.targetMessage;

  if (!isAdmin) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("権限がありません")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [errorEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  if (
    !interaction.channel ||
    !interaction.channel.isTextBased() ||
    interaction.channel instanceof PartialGroupDMChannel
  ) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("このチャンネルではメッセージをピン留めできません。")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [errorEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  if (targetMsg.author?.bot) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription("ボットのメッセージはピン留めできません。")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [errorEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const dataManager = new ServerDataManager(interaction.guildId as string);
  const existingPinned = await dataManager.getConfig("pinnedMessage", interaction.channelId);
  if (existingPinned) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("エラー")
      .setDescription(
        "このチャンネルには既にピン留めされたメッセージがあります。\n最初にそれを`/unpin`で解除してください。"
      )
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({
      embeds: [errorEmbed],
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  const embed = new EmbedBuilder()
    .setDescription(interaction.targetMessage.content)
    .setColor(config.color.success)
    .setFooter({ text: "Pinned Message" });
  const sendedMessage = await interaction.channel.send({ embeds: [embed] });
  dataManager.setConfig(
    "pinnedMessage",
    { message: interaction.targetMessage.content, latestMessageId: sendedMessage.id },
    interaction.channelId
  );

  const successEmbed = new EmbedBuilder()
    .setTitle("メッセージをピン留めしました")
    .setDescription(
      "このメッセージは今後ピン留めされます。\nファイルは保存されないのでご注意ください。"
    )
    .setColor(config.color.success)
    .setTimestamp();

  await interaction.reply({
    embeds: [successEmbed],
    flags: [MessageFlags.Ephemeral],
  });
};
