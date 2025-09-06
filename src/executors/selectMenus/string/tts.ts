import { EmbedBuilder, MessageFlags, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";
import { subSelectMenusHandling } from "@/lib/commandUtils";
import { GetErrorChannel } from "@/lib/channelUtils";
import { ALStorage, loggingSystem } from "@/index";

export const name = "tts";

const handlingCommands = subSelectMenusHandling("string/tts");

export const execute = async (interaction: StringSelectMenuInteraction) => {
  const ctx = ALStorage.getStore();
  const logger = loggingSystem.getLogger({ ...ctx, function: "selectMenus/string/tts" });
  const commandKey = interaction.customId.split("_")[1];
  const commands = handlingCommands.get(commandKey);

  if (!commands) {
    logger.error(
      { service: "TTS", userId: interaction.user.id, commandKey },
      "No command found for TTS select menu"
    );
    return;
  }

  try {
    await commands.execute(interaction);
  } catch (error) {
    logger.error({ error, stack_trace: (error as Error).stack }, `An Error Occurred`, error);

    const logEmbed = new EmbedBuilder()
      .setTitle("ERROR - cmd")
      .setDescription("```\n" + String(error) + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    const errorChannel = await GetErrorChannel(interaction.client);
    if (errorChannel) {
      errorChannel.send({ embeds: [logEmbed] });
    }

    const userEmbed = new EmbedBuilder()
      .setTitle("すみません。エラーが発生しました。")
      .setDescription("```\n" + String(error) + "\n```")
      .setColor(config.color.error)
      .setTimestamp();

    if (interaction.channel && interaction.channel.isSendable()) {
      await interaction.channel.send({
        embeds: [userEmbed],
        flags: MessageFlags.SuppressEmbeds,
      });
    }
  }
};
