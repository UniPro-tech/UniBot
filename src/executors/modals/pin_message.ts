import config from "@/config";
import { ServerDataManager } from "@/lib/dataUtils";
import {
  EmbedBuilder,
  MessageFlags,
  ModalSubmitInteraction,
  PartialGroupDMChannel,
} from "discord.js";

export const name = "pin_message";

export const execute = async (interaction: ModalSubmitInteraction) => {
  const message = interaction.fields.getTextInputValue("message");
  const channelId = interaction.channelId as string;

  const channel = interaction.channel;
  if (channel && channel.isTextBased() && channel instanceof PartialGroupDMChannel === false) {
    const embed = new EmbedBuilder()
      .setDescription(message)
      .setColor(config.color.success)
      .setFooter({ text: "Pinned Message" });
    const sendedMessage = await channel.send({ embeds: [embed] });
    const dataManager = new ServerDataManager(interaction.guildId as string);
    dataManager.setConfig(
      "pinnedMessage",
      { message, latestMessageId: sendedMessage.id },
      channelId
    );
  } else {
    await interaction.reply({
      content: `このチャンネルではメッセージを送信できません。`,
      flags: [MessageFlags.Ephemeral],
    });
    return;
  }

  await interaction.reply({
    content: `メッセージをピン留めしました: \`${message}\``,
    flags: [MessageFlags.Ephemeral],
  });
};
