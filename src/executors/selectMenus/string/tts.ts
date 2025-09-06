import { EmbedBuilder, MessageFlags, StringSelectMenuInteraction } from "discord.js";
import config from "@/config";
import { subSelectMenusHandling } from "@/lib/commandUtils";
import { GetErrorChannel } from "@/lib/channelUtils";
import { logger } from "@/index";

export const name = "tts";

const handlingCommands = subSelectMenusHandling("string/tts");

export const execute = async (interaction: StringSelectMenuInteraction) => {
  const [, commandKey] = interaction.customId.split("_");
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
    const time = interaction.client.functions.timeUtils.timeToJSTstamp(Date.now(), true);
    logger.error(
      { service: "TTS", userId: interaction.user.id, commandKey, time, error },
      `An Error Occurred in ${interaction.customId}`
    );

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
