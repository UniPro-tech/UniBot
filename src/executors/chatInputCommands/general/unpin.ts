import config from "@/config";
import { ServerDataManager } from "@/lib/dataUtils";
import { EmbedBuilder, SlashCommandBuilder } from "@discordjs/builders";
import {
  ChatInputCommandInteraction,
  InteractionContextType,
  PermissionFlagsBits,
} from "discord.js";

export const guildOnly = false;

export const data = new SlashCommandBuilder()
  .setName("unpin")
  .setDescription("ピン留めを解除します。")
  .setDefaultMemberPermissions(PermissionFlagsBits.PinMessages)
  .setContexts(InteractionContextType.Guild);

export const execute = async (interaction: ChatInputCommandInteraction) => {
  try {
    const dataManager = new ServerDataManager(interaction.guildId!);
    const channelId = interaction.channelId;
    const pinnedMessageConfig = await dataManager.getConfig("pinnedMessage", channelId);
    if (!pinnedMessageConfig) {
      await interaction.reply({
        content: "このチャンネルにはピン留めされたメッセージがありません。",
        ephemeral: true,
      });
      return;
    }

    await dataManager.deleteConfig("pinnedMessage", channelId);

    const successEmbed = new EmbedBuilder()
      .setTitle("ピン留めを解除しました")
      .setColor(config.color.success)
      .setTimestamp();

    await interaction.reply({ embeds: [successEmbed], ephemeral: true });
  } catch (error) {
    const errorEmbed = new EmbedBuilder()
      .setTitle("エラーが発生しました")
      .setDescription("ピン留めの解除中にエラーが発生しました。")
      .setColor(config.color.error)
      .setTimestamp();
    await interaction.reply({ embeds: [errorEmbed], ephemeral: true });
  }
};

export default {
  guildOnly,
  data,
  execute,
};
